"use strict";(globalThis["webpackChunkcontrol"]=globalThis["webpackChunkcontrol"]||[]).push([[111],{2111:(e,n,t)=>{t.r(n),t.d(n,{default:()=>S});var a=t(9835),o=t(499),l=t(1957),s=t(4304),u=t(2626),r=t(8339),i=t(6647),c=t(9036),d=function(e,n,t,a){function o(e){return e instanceof t?e:new t((function(n){n(e)}))}return new(t||(t=Promise))((function(t,l){function s(e){try{r(a.next(e))}catch(n){l(n)}}function u(e){try{r(a["throw"](e))}catch(n){l(n)}}function r(e){e.done?t(e.value):o(e.value).then(s,u)}r((a=a.apply(e,n||[])).next())}))};class g extends c.Z{authentication(e){return d(this,void 0,void 0,(function*(){return this.timeoutInMilliseconds=1e4,this.request.post("v1/auth/login/",e)}))}}var p=t(503),v=t(3578),h=t(8900),f=t(9302);const m={class:"flex justify-center"},w={class:"q-gutter-y-md q-mt-xs"},b=(0,a.aZ)({__name:"LoginIndex",setup(e){const n=new p.Z,t=new v.Z,c=(0,r.tv)(),d=(0,i.QT)().t,b=(0,o.iH)(""),x=(0,o.iH)(""),y=(0,o.iH)(!1),I=(0,a.Fl)((()=>b.value.length>=3&&x.value.length>=6));function Z(){n.clearAllLocalStorage()}function k(e){t.setToken(e)}function C(e){n.setAccountId(e)}function _(e){n.setUsername(e)}function q(e){const n=e.split(".")[1],t=n.replace("-","+").replace("_","/");return JSON.parse(atob(t))}const S=(0,a.Fl)((()=>{const e=(0,f.Z)();return e.dark.isActive?"/_/assets/control-logo-dark.svg":"/_/assets/control-logo-light.svg"}));function U(){y.value=!0;const e=new g;y.value=!0,e.authentication({username:b.value,password:x.value}).then((e=>{const n=e.data.body.tokenStr,a=q(n);Z(),t.removeToken(),k(e.data.body.tokenStr),C(a.accountId),_(b.value),(0,h.LX)(d("loginIndex.messageLoginSuccess")),c.push("/overview")})).catch((e=>{console.error(e),(0,h.s9)(e.response.data,d("loginIndex.messageLoginError"))})).finally((()=>{y.value=!1}))}return(e,n)=>{const t=(0,a.up)("q-img"),r=(0,a.up)("q-btn"),i=(0,a.up)("q-card-section"),c=(0,a.up)("q-card"),g=(0,a.up)("q-page-container");return(0,a.wg)(),(0,a.j4)(g,null,{default:(0,a.w5)((()=>[(0,a.Wm)(c,{flat:"",class:"absolute-center no-shadow login-bg-blur login-card"},{default:(0,a.w5)((()=>[(0,a._)("div",m,[(0,a.Wm)(t,{src:S.value,style:{height:"50px",width:"250px",top:"-25px"}},null,8,["src"])]),(0,a.Wm)(i,null,{default:(0,a.w5)((()=>[(0,a._)("div",w,[(0,a.Wm)(s.Z,{value:b.value,"onUpdate:value":n[0]||(n[0]=e=>b.value=e),label:(0,o.SU)(d)("loginIndex.usernameInput"),dataCyInput:"login-username",icon:"person",rules:[()=>""!==b.value||e.$t("accountsCreateUserDialog.usernameRequired"),()=>b.value.length>=3||e.$t("accountsCreateUserDialog.usernameMinLength")]},null,8,["value","label","rules"]),(0,a.Wm)(u.Z,{dataCyInput:"login-password",label:(0,o.SU)(d)("loginIndex.passwordInput"),password:x.value,"onUpdate:password":n[1]||(n[1]=e=>x.value=e),onKeyup:n[2]||(n[2]=(0,l.D2)((e=>U()),["enter"]))},null,8,["label","password"]),(0,a.Wm)(r,{loading:y.value,disable:!1===I.value,icon:"login",color:!1===I.value?"grey-8":"positive",class:"full-width",size:"lg",onClick:n[3]||(n[3]=e=>U())},null,8,["loading","disable","color"])])])),_:1})])),_:1})])),_:1})}}});var x=t(2133),y=t(4458),I=t(335),Z=t(3190),k=t(8879),C=t(9984),_=t.n(C);const q=b,S=q;_()(b,"components",{QPageContainer:x.Z,QCard:y.Z,QImg:I.Z,QCardSection:Z.Z,QBtn:k.Z})}}]);