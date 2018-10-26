/*
 * Copyright 1999-2018 Alibaba Group.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package downloader

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/alibaba/Dragonfly/dfget/config"
	"github.com/alibaba/Dragonfly/dfget/core/api"
	"github.com/alibaba/Dragonfly/dfget/core/regist"
	"github.com/alibaba/Dragonfly/dfget/types"
	"github.com/alibaba/Dragonfly/dfget/util"
)

// P2PDownloader is one implementation of Downloader that uses p2p pattern
// to download files.
type P2PDownloader struct {
	Ctx            *config.Context
	API            api.SupernodeAPI
	Register       regist.SupernodeRegister
	RegisterResult *regist.RegisterResult

	node         string
	taskID       string
	targetFile   string
	taskFileName string

	pieceSizeHistory [2]int32
	queue            util.Queue
	clientQueue      util.Queue
	writerDone       chan struct{}

	// pieceSet range -> bool
	// true: if the range is processed successfully
	// false: if the range is in processing
	// not in: the range hasn't been processed
	pieceSet map[string]bool
	total    int64
}

func (p2p *P2PDownloader) init() {
	p2p.node = p2p.RegisterResult.Node
	p2p.taskID = p2p.RegisterResult.TaskID
	p2p.targetFile = p2p.Ctx.RV.RealTarget
	p2p.taskFileName = p2p.Ctx.RV.TaskFileName

	p2p.pieceSizeHistory[0], p2p.pieceSizeHistory[1] =
		p2p.RegisterResult.PieceSize, p2p.RegisterResult.PieceSize

	p2p.queue = util.NewQueue(0)
	p2p.clientQueue = util.NewQueue(config.DefaultClientQueueSize)
	p2p.writerDone = make(chan struct{})

	p2p.pieceSet = make(map[string]bool)
}

// Run starts to download the file.
func (p2p *P2PDownloader) Run() error {
	var (
		lastItem *PieceItem
		goNext   bool
	)

	for {
		goNext, lastItem = p2p.getItem(lastItem)
		if !goNext {
			continue
		}
		p2p.Ctx.ClientLogger.Infof("p2p download:%v", lastItem)

		curItem := *lastItem
		curItem.PieceContents = bytes.Buffer{}
		lastItem = nil

		response, err := p2p.pullPieceTask(&curItem)
		if err == nil {
			code := response.Code
			if code == config.TaskCodeContinue {
				p2p.processPiece(response, &curItem)
			} else if code == config.TaskCodeFinish {
				p2p.finishTask(response)
				return nil
			} else {
				p2p.Ctx.ClientLogger.Warnf("request piece result:%v", response)
				if code == config.TaskCodeSourceError {
					p2p.Ctx.BackSourceReason = config.BackSourceReasonSourceError
				}
			}
		} else {
			p2p.Ctx.ClientLogger.Errorf("p2p download fail: %v", err)
			if p2p.Ctx.BackSourceReason == 0 {
				p2p.Ctx.BackSourceReason = config.BackSourceReasonDownloadError
			}
		}

		if p2p.Ctx.BackSourceReason != 0 {
			backDownloader := NewBackDownloader(p2p.Ctx, p2p.RegisterResult)
			return backDownloader.Run()
		}
	}
}

// Cleanup clean all temporary resources generated by executing Run.
func (p2p *P2PDownloader) Cleanup() {
}

func (p2p *P2PDownloader) pullPieceTask(item *PieceItem) (
	*types.PullPieceTaskResponse, error) {
	var (
		res *types.PullPieceTaskResponse
		err error
	)
	req := &types.PullPieceTaskRequest{
		SrcCid: p2p.Ctx.RV.Cid,
		DstCid: item.DstCid,
		Range:  item.Range,
		Result: item.Result,
		Status: item.Status,
		TaskID: item.TaskID,
	}

	for {
		if res, err = p2p.API.PullPieceTask(item.SuperNode, req); err != nil {
			p2p.Ctx.ClientLogger.Errorf("pull piece task error: %v", err)
		} else if res.Code == config.TaskCodeWait {
			sleepTime := time.Duration(rand.Intn(1400)+600) * time.Millisecond
			p2p.Ctx.ClientLogger.Infof("pull piece task result:%s and sleep %.3fs",
				res, sleepTime.Seconds())
			time.Sleep(sleepTime)
			continue
		}
		break
	}

	if res == nil || (res.Code != config.TaskCodeContinue &&
		res.Code != config.TaskCodeFinish &&
		res.Code != config.TaskCodeLimited &&
		res.Code != config.Success) {
		p2p.Ctx.ClientLogger.Errorf("pull piece task fail:%v and will migrate", res)

		var registerRes *regist.RegisterResult
		if registerRes, err = p2p.Register.Register(p2p.Ctx.RV.PeerPort); err != nil {
			return nil, err
		}
		p2p.pieceSizeHistory[1] = registerRes.PieceSize
		item.Status = config.TaskStatusStart
		item.SuperNode = registerRes.Node
		item.TaskID = registerRes.TaskID
		util.Printer.Println("migrated to node:" + item.SuperNode)
		return p2p.pullPieceTask(item)
	}

	return res, err
}

func (p2p *P2PDownloader) pullRate(data *types.PullPieceTaskResponseContinueData) {

}

func (p2p *P2PDownloader) startTask(data *types.PullPieceTaskResponseContinueData) {
}

func (p2p *P2PDownloader) getItem(latestItem *PieceItem) (bool, *PieceItem) {
	var (
		needMerge = true
	)
	if ok, v := p2p.queue.PollTimeout(2 * time.Second); ok {
		item := v.(*PieceItem)
		if item.PieceSize != 0 && item.PieceSize != p2p.pieceSizeHistory[1] {
			return false, latestItem
		}
		if item.SuperNode != p2p.node {
			item.DstCid = ""
			item.SuperNode = p2p.node
			item.TaskID = p2p.taskID
		}
		if item.Range != "" {
			ok, v := p2p.pieceSet[item.Range]
			if !ok {
				p2p.Ctx.ClientLogger.Warnf("pieceRange:%s is neither running nor success")
				return false, latestItem
			}
			if !v && (item.Result == config.ResultSemiSuc ||
				item.Result == config.ResultSuc) {
				p2p.total += int64(item.PieceContents.Len())
				p2p.pieceSet[item.Range] = true
			} else if !v {
				delete(p2p.pieceSet, item.Range)
			}
		}
		latestItem = item
	} else {
		p2p.Ctx.ClientLogger.Warnf("get item timeout(2s) from queue.")
		needMerge = false
	}
	if util.IsNil(latestItem) {
		return false, latestItem
	}
	if latestItem.Result == config.ResultSuc ||
		latestItem.Result == config.ResultFail ||
		latestItem.Result == config.ResultInvalid {
		needMerge = false
	}
	runningCount := 0
	for _, v := range p2p.pieceSet {
		if !v {
			runningCount++
		}
	}
	if needMerge && (p2p.queue.Len() > 0 || runningCount > 2) {
		return false, latestItem
	}
	return true, latestItem
}

func (p2p *P2PDownloader) processPiece(response *types.PullPieceTaskResponse,
	item *PieceItem) {

}

func (p2p *P2PDownloader) finishTask(response *types.PullPieceTaskResponse) {

}

func (p2p *P2PDownloader) refresh(item *PieceItem) {

}
