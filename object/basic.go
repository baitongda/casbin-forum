// Copyright 2020 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dchest/captcha"

	"github.com/casbin/casbin-forum/util"
)

type BasicInfo struct {
	Id    string `xorm:"varchar(100) notnull pk"`
	Value string `xorm:"mediumtext"`
}

var fileDate, version string

func GetForumVersion() string {
	pwd, _ := os.Getwd()

	fileInfos, err := ioutil.ReadDir(pwd + "/.git/refs/heads")
	for _, v := range fileInfos {
		if v.Name() == "master" {
			if v.ModTime().String() == fileDate {
				return version
			} else {
				fileDate = v.ModTime().String()
				break
			}
		}
	}

	content, err := ioutil.ReadFile(pwd + "/.git/refs/heads/master")
	if err != nil {
		return ""
	}

	//Convert to full length
	temp := string(content)
	version = strings.ReplaceAll(temp, "\n", "")

	return version
}

func GetHighestOnlineNum() int {
	info := BasicInfo{Id: "HighestOnlineNum"}
	existed, err := adapter.engine.Get(&info)
	if err != nil {
		panic(err)
	}

	if existed {
		return util.ParseInt(info.Value)
	} else {
		info := BasicInfo{
			Id:    "HighestOnlineNum",
			Value: "0",
		}

		_, err := adapter.engine.Insert(&info)
		if err != nil {
			panic(err)
		}

		return 0
	}
}

func UpdateHighestOnlineNum(num int) bool {
	info := new(BasicInfo)
	info.Value = util.IntToString(num)
	affected, err := adapter.engine.Where("id = ?", "HighestOnlineNum").Cols("value").Update(info)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetCaptcha() (string, []byte) {
	id := captcha.NewLen(5)

	var buffer bytes.Buffer

	err := captcha.WriteImage(&buffer, id, 320, 80)
	if err != nil {
		panic(err)
	}

	return id, buffer.Bytes()
}

func VerifyCaptcha(id, digits string) bool {
	res := captcha.VerifyString(id, digits)

	return res
}

func GetCronJobs() []*CronJob {
	info := BasicInfo{Id: "CronJobs"}
	existed, err := adapter.engine.Get(&info)
	if err != nil {
		panic(err)
	}

	if existed {
		var jobs []*CronJob
		err := json.Unmarshal([]byte(info.Value), &jobs)
		if err != nil {
			panic(err)
		}
		return jobs
	} else {
		jobs, err := json.Marshal(DefaultCronJobs)
		if err != nil {
			panic(err)
		}
		info := BasicInfo{
			Id:    "CronJobs",
			Value: string(jobs),
		}

		_, err = adapter.engine.Insert(&info)
		if err != nil {
			panic(err)
		}

		return DefaultCronJobs
	}
}

func GetCronPosts() []*PostJob {
	info := BasicInfo{Id: "CronPosts"}
	existed, err := adapter.engine.Get(&info)
	if err != nil {
		panic(err)
	}

	if existed {
		var posts []*PostJob
		err := json.Unmarshal([]byte(info.Value), &posts)
		if err != nil {
			panic(err)
		}
		return posts
	} else {
		posts, err := json.Marshal(DefaultCronPosts)
		if err != nil {
			panic(err)
		}
		info := BasicInfo{
			Id:    "CronPosts",
			Value: string(posts),
		}

		_, err = adapter.engine.Insert(&info)
		if err != nil {
			panic(err)
		}

		return DefaultCronPosts
	}
}