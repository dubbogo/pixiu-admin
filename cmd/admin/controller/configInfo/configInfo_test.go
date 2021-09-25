/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package configInfo

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_method1(t *testing.T) {
	var list1, list2 []string
	list1 = append(list1, "a")
	list2 = nil
	for _, v1 := range list1 {
		for _, v := range list2 {
			fmt.Println("in side" + v)
			//t.Log("inside" + v)
		}
		fmt.Println(v1)
		//t.Log("outside" + v1)
	}
}

func Test_regx_split(t *testing.T) {
	txt := "/config/api/resources/1/xxx"
	re := regexp.MustCompile("^([^s]*)/resources/")
	split := re.Split(txt, -1)
	fmt.Println(split)
}
