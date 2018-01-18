package utils

import (
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

//IdWorker 唯一编号生成器
var IdWorker *IdWorkerStruct

func init() {
	IdWorker = idWorker(5)
}

//唯一序号生成对象
type IdWorkerStruct struct {
	workerId           int    //机器ID
	twepoch            int64  //唯一时间，这是一个避免重复的随机量，自行设定不要大于当前时间戳
	sequence           uint32 //当前时间戳计数器
	workerIdBits       uint   //机器码字节数,默认4字节保存机器ID
	maxWorkerId        int    //最大机器ID
	sequenceBits       uint   //计数器字节数,默认10个字节保存计数器
	workerIdShift      uint   //机器码数据左移位数
	timestampLeftShift uint   //时间戳左移位数
	sequenceMask       uint   //一微秒内可以产生的计数器值,达到该值后要等到下一微妙再生成
	lastTimestamp      int64  //上次生成序号的时间戳
}

//生成一个IdWOkerStruct对象
func idWorker(workerId int, params ...int64) *IdWorkerStruct {
	var sequenceBits uint
	var twepoch int64
	if len(params) > 1 {
		sequenceBits = uint(params[0])
		twepoch = int64(params[1])
	} else {
		if len(params) > 0 {
			sequenceBits = uint(params[0])
			twepoch = 1359338079273673920
		} else {
			sequenceBits = 1
			twepoch = 1359338079273673920
		}
	}
	idw := new(IdWorkerStruct)
	idw.workerId = workerId
	idw.twepoch = twepoch
	idw.sequenceBits = sequenceBits
	idw.sequence = 0
	idw.workerIdBits = 10
	idw.maxWorkerId = -1 << idw.workerIdBits
	tp := -1 << idw.sequenceBits
	idw.sequenceMask = uint(tp)
	idw.workerIdShift = idw.sequenceBits
	idw.timestampLeftShift = idw.workerIdBits + idw.workerIdShift
	idw.lastTimestamp = -1
	return idw
}

//生成一个唯一ID
func (d *IdWorkerStruct) Next() (uint64, error) {
	timestamp := timeGen()
	if d.lastTimestamp == timestamp {
		newSequenceMask := atomic.AddUint32(&d.sequence, 1)
		isend := uint(newSequenceMask) & d.sequenceMask
		if isend == 0 {
			timestamp = tilNextMillis(d.lastTimestamp)
		}
	} else {
		d.sequence = 0
	}
	if timestamp < d.lastTimestamp {
		return 0, errors.New(fmt.Sprint("Clock moved backwards.  Refusing to generate id for %d milliseconds", d.lastTimestamp))
	}
	d.lastTimestamp = timestamp
	i := ((timestamp - d.twepoch) << d.timestampLeftShift) | int64((d.workerId << d.workerIdShift)) | int64(d.sequence)
	return uint64(i), nil
}

func (d *IdWorkerStruct) NextString() (string, error) {
	id, err := d.Next()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}

//获取下一微秒时间戳
func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for {
		if timestamp > lastTimestamp {
			break
		}
		timestamp = timeGen()
	}
	return timestamp
}

//当前时间戳
func timeGen() int64 {
	return time.Now().UnixNano()
}
