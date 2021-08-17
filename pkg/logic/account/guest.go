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

package account

import (
	"github.com/dubbogo/pixiu-admin/pkg/dao"
	"github.com/dubbogo/pixiu-admin/pkg/dao/impl"
)

var(
	guestDao dao.GuestDao = impl.NewGuestDao()
	userDao dao.UserDao = impl.NewUserDao()
)

func Login(username string, password string) (bool, int) {
	return guestDao.Login(username, password)
}

func Register(username, password string) error {
	return guestDao.Register(username, password)
}
