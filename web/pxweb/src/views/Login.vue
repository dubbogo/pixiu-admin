<template>
    <div class="login-container" id="container"> 
        <div >
            <div class="login-header"></div>
            <div class="card-box">
                <h1>pixiu控制台管理系统</h1>
                <el-form
                        class="login-form"
                        autoComplete="on"
                        :model="ruleForm"
                        :rules="rules"
                        ref="ruleForm"
                        label-position="left">
                    <div class="login_msg" v-if="loginMsg">
                        <i class="el-icon-error"></i>
                        <p class="error">
                            {{loginMsg}}
                        </p>
                    </div>
                    <el-form-item prop="userName" class="item">
                        <el-input
                                placeholder="请输入用户名"
                                name="userName"
                                autoComplete="on"
                                v-model="ruleForm.userName">
                            <i slot="prefix" class="el-icon-etcyonghuming"></i>
                        </el-input>
                    </el-form-item>
                    <el-form-item
                            prop="password"
                            class="item">
                        <el-input
                                placeholder="请输入密码"
                                name="password"
                                :type="isShowPwd ? 'text' : 'password'"
                                @keyup.enter.native="handleLogin"
                                v-model="ruleForm.password"
                                autoComplete="on">
                            <i slot="prefix" class="el-icon-etcyonghuming1"></i>
                        </el-input>
                    </el-form-item>
                    <div class="vCode" v-if="code.isVCode">
                        <el-form-item
                                prop="verification_code"
                                class="item">
                            <el-input
                                    placeholder="请输入验证码"
                                    name="app_ver"
                                    type="text"
                                    v-model="ruleForm.verification_code">
                                <i slot="prefix" class="el-icon-etcyanzhengma"></i>
                            </el-input>
                        </el-form-item>
                        <div id="code" title="看不清？换一张～" @click="handleOnRandom">
                            <img :src="'/issue/captcha?n='+ rand"/>
                        </div>
                    </div>

                    <div>
                        <el-button type="primary" style="width:100%;margin-bottom:30px;" :loading="loading"
                                @click.native="handleLogin()">登录
                        </el-button>
                    </div>
                </el-form>
            </div>
            <div class="login-footer">
                技术支持
            </div>
        </div>
        
    </div>
</template>

<script>
    import api from '@/api'
    import {mapActions} from 'vuex';
    import {loading} from "../utils/dialogUtils";
    import {getSecuCode} from "../utils/dialogUtils";
import { getToken, setToken, removeToken, getLocalStorage, setLocalStorage,clearLocalStorage } from '@/utils/auth'
    import THREE from "@/utils/three";
    var SEPARATION = 100, AMOUNTX = 60, AMOUNTY = 40;
    var container;
    var camera, scene, renderer;
    var particles, particle, count = 0;
    var mouseX = 0, mouseY = 0;
    var windowHalfX = window.innerWidth / 2;
    var windowHalfY = window.innerHeight / 2;
    export default {
        data() {
            let validatePwd = (rule, value, callback) => {

                if (value === "" || value === undefined) {
                    callback(new Error("请输入密码"));
                } else {
                    callback();
                }
            };
            return {
                ruleForm: {
                    userName: "admin",
                    password: "admin",
                    checked: true
                },
                rules: {
                    userName: [
                        {required: true, message: "请输入登录名", trigger: "blur"}
                    ],
                    password: [{validator: validatePwd, trigger: "blur"}]
                },
                code: {
                    isVCode: false,//是否显示验证码
                    src: ''
                },
                rand: 0,
                error12: 0,
                loginMsg: '',
                isShowPwd: false, // 是否显示密码
                loading: false, // 登录loading
            };
        },
        methods: {
            ...mapActions([
                'Login'
            ]),
             init() {

				// container = document.createElement( 'div' );	//创建容器
				// document.body.appendChild( container );			//将容器添加到页面上
                container = document.getElementById( 'container' );	//创建容器
				camera = new THREE.THREE.PerspectiveCamera( 120, window.innerWidth / window.innerHeight, 1, 1500 );		//创建透视相机设置相机角度大小等
				camera.position.set(0,450,2000);		//设置相机位置

				scene = new THREE.THREE.Scene();			//创建场景
				particles = new Array();

				var PI2 = Math.PI * 2;
				//设置粒子的大小，颜色位置等
				var material = new THREE.THREE.ParticleCanvasMaterial( {
					color: 0x0f96ff,
					vertexColors:true,
					size: 4,
					program: function ( context ) {
						context.beginPath();
						context.arc( 0, 0, 0.01, 0, PI2, true );	//画一个圆形。此处可修改大小。
						context.fill();
					}
				} );
				//设置长条粒子的大小颜色长度等
				var materialY = new THREE.THREE.ParticleCanvasMaterial( {
					color: 0xffffff,
					vertexColors:true,
					size: 1,
					program: function ( context ) {

						context.beginPath();
						//绘制渐变色的矩形
						var lGrd = context.createLinearGradient(-0.008,0.25,0.016,-0.25);
						lGrd.addColorStop(0, '#16eff7');
						lGrd.addColorStop(1, '#0090ff');
						context.fillStyle = lGrd;
						context.fillRect(-0.008,0.25,0.016,-0.25);  //注意此处的坐标大小
						//绘制底部和顶部圆圈
						context.fillStyle = "#0090ff";
						context.arc(0, 0, 0.008, 0, PI2, true);    //绘制底部圆圈
						context.fillStyle = "#16eff7";
						context.arc(0, 0.25, 0.008, 0, PI2, true);    //绘制顶部圆圈
						context.fill();
						context.closePath();
						//绘制顶部渐变色光圈
						var rGrd = context.createRadialGradient(0, 0.25, 0, 0, 0.25, 0.025);
						rGrd.addColorStop(0, 'transparent');
						rGrd.addColorStop(1, '#16eff7');
						context.fillStyle = rGrd;
						context.arc(0, 0.25, 0.025, 0, PI2, true);    //绘制一个圆圈
						context.fill();

					}
				} );

				//循环判断创建随机数选择创建粒子或者粒子条
				var i = 0;
				for ( var ix = 0; ix < AMOUNTX; ix ++ ) {
					for ( var iy = 0; iy < AMOUNTY; iy ++ ) {
						var num = Math.random()-0.1;
						if (num >0 ) {
							particle = particles[ i ++ ] = new THREE.THREE.Particle( material );
							// console.log("material")
						}
						else {
							particle = particles[ i ++ ] = new THREE.THREE.Particle( materialY );
							// console.log("materialY")
						}
						//particle = particles[ i ++ ] = new THREE.Particle( material );
						particle.position.x = ix * SEPARATION - ( ( AMOUNTX * SEPARATION ) / 2 );
						particle.position.z = iy * SEPARATION - ( ( AMOUNTY * SEPARATION ) / 2 );
						scene.add( particle );
					}
				}

				renderer = new THREE.THREE.CanvasRenderer();
				renderer.setSize( window.innerWidth, window.innerHeight );
				container.appendChild( renderer.domElement );
				//document.addEventListener( 'mousemove', onDocumentMouseMove, false );
				//document.addEventListener( 'touchstart', onDocumentTouchStart, false );
				//document.addEventListener( 'touchmove', onDocumentTouchMove, false );
				window.addEventListener( 'resize', this.onWindowResize, false );
			},

			//浏览器大小改变时重新渲染
			onWindowResize() {
				windowHalfX = window.innerWidth / 2;
				windowHalfY = window.innerHeight / 2;
				camera.aspect = window.innerWidth / window.innerHeight;
				camera.updateProjectionMatrix();
				renderer.setSize( window.innerWidth, window.innerHeight );
			},
             animate() {
				requestAnimationFrame( this.animate );
				this.render();
			},

			//将相机和场景渲染到页面上
			render() {
				var i = 0;
				//更新粒子的位置和大小
				for (var ix = 0; ix < AMOUNTX; ix++) {
					for (var iy = 0; iy < AMOUNTY; iy++) {
						particle = particles[i++];
						//更新粒子位置
						particle.position.y = (Math.sin((ix + count) * 0.3) * 50) + (Math.sin((iy + count) * 0.5) * 50);
						//更新粒子大小
						particle.scale.x =  particle.scale.y = particle.scale.z  = ( (Math.sin((ix + count) * 0.3) + 1) * 4 + (Math.sin((iy + count) * 0.5) + 1) * 4 )*50;	//正常情况下再放大100倍*1200
					}
				}

				renderer.render( scene, camera );
				count += 0.1;
			},
            handleLogin() {
                let data = {
                    code: 200,
                    msg:'登录成功',
                    data:{
                        userName:'admin'
                    }
                }
                setToken('111')
                setLocalStorage('operatorInfo',data);
                // this.$router.push('/dashBoard/manage');
                this.$router.push({
                    path: '/manage',
                    query: {
                        // id: new Date().getTime()
                    }
                })
                // return;
                // this.$refs["ruleForm"].validate(valid => {
                //     if (valid) {
                //         let params = {
                //             userName: this.ruleForm.userName,
                //             userPwd: this.ruleForm.password,
                //         }

                //         this.Login(params).then(res => {
                //             this.loginMsg = '';
                //              console.log('res'+res)
                //            if(getLocalStorage('operatorInfo').code==200){
                //                 this.$router.push({name: 'DashBoard'});
                //            }else if(getLocalStorage('operatorInfo').code!=200){
                //                 this.loginMsg = getLocalStorage('operatorInfo').msg;
                //                 return;
                //            }
                //         }, error => {
                //             console.log(error)
                //             this.loginMsg = error.msg;
                //         });
                //     } else {
                //         return false;
                //     }
                // });
            },
            handleOnRandom() {
                this.rand += 1;
            }
        },
        mounted() {
            this.init();		//初始化
			this.animate();	//动画效果
        }
    };
</script>

<style type="text/scss" lang="scss">
    @import "../styles/mixin";
    @import "../styles/common";


    $dark_gray: #889aa4;
    $light_gray: #eee;

    .login-container {
        @include relative;
        margin: 0;
        overflow: hidden;
        background: linear-gradient(to bottom, #19778c, #095f88);
        background-size:1400% 300%;
        animation: dynamics 6s ease infinite;
        -webkit-animation: dynamics 6s ease infinite;
        -moz-animation: dynamics 6s ease infinite;
        font-size: 14px;
        color: #ffffff;
        min-height: 700px;
        height: 100%;

        .login-header {
            width: 100%;
            height: 50%;
            // background: url('../assets/login_bg.png') no-repeat center #fff;
            background-size: 100% 100%;
        }

        .card-box {
            @include fxied-center;
            text-align: center;
            width: 60%;
            transform: translate(-50%, -57%);

            h1 {
                color: #f5f5f5;
                margin-bottom: 20px;
                text-shadow: 2px 2px 5px #0c0c0c;
            }

            .login-form {
                padding: 51px 48px 48px 48px;
                background: rgba(255, 255, 255, 0.2);
                border-radius: 5px;
                box-shadow: 1px 1px 3px #999;
                width: 35%;
                margin: 0 auto;
            }

            .login_msg {
                background-color: #fef2f2;
                color: #6C6C6C;
                line-height: 16px;
                padding: 6px 10px;
                background: #fef2f2;
                border: 1px solid #ffb4a8;
                margin-bottom: 22px;
                overflow: hidden;
                text-align: left;
                display: flex;

                i {

                    padding-right: 10px;
                    color: #f40;
                }

                p {
                    @include f_left;
                    white-space: normal;
                    word-wrap: break-word;
                    padding: 0;
                    margin: 0;
                }
            }
        }

        .item {
            .el-form-item__content {
                display: flex;
                flex-flow: row;

                .el-input__inner {
                    color: #999;
                }
            }

        }

        .vCode {
            display: flex;
            flex-flow: row;
            justify-content: space-between;

            .item {
                width: 60%;
            }

            #code {
                display: block;
                color: #ffffff;
                font-size: 20px;
                padding: 5px 35px 10px 35px;
                margin-left: 5%;
                height: 27px;
                cursor: pointer;
            }
        }

        input {
            border: 0;
            -webkit-appearance: none;
            color: $light_gray;
            height: 100%;
        }

        .el-input {
            display: inline-block;
        }

        .tips {
            font-size: 14px;
            color: #fff;
            margin-bottom: 0.13333rem;
        }

        .svg-container {
            padding: 0.08rem 0.0666rem 0.08rem 0.2rem;
            color: $dark_gray;
            vertical-align: middle;
            display: inline-block;

            &_login {
                font-size: 20px;
            }
        }

        .title {
            font-size: 26px;
            color: $light_gray;
            margin: 0 auto 0.5333rem auto;
            text-align: center;
            font-weight: bold;
        }

        .el-form-item {
            border: 1px solid #cbcbcb;
            background: #fff;
            border-radius: 5px;
            color: #9f9f9f;
        }

        .show-pwd {
            position: absolute;
            right: 0.1333rem;
            top: 0.09333rem;
            font-size: 16px;
            color: $dark_gray;
            cursor: pointer;
        }

        .login-footer {
            position: absolute;
            left: 0;
            right: 0;
            bottom: 0;
            height: 32px;
            text-align: center;
            color: #143992;
            // background: url('../assets/login-ft.png') no-repeat center;
        }

        .SoftwareDownload{
            display: flex;
            flex-direction: row;
            justify-content: flex-start;
            .SoftwareDownloadTitle{
                margin-right: 10px;
            }
            .SoftwareDownloadClick{
                flex: 1;
                display: flex;
                flex-direction: row;
                justify-content: flex-start;

            }
            span{
                color: $three-color;
                margin-right: 8px;
                cursor: pointer;
            }
            .line{
                font-weight: 500;
                color: black;
            }
        }
    }
</style>
        