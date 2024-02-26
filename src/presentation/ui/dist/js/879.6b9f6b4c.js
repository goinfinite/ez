"use strict";(globalThis["webpackChunkcontrol"]=globalThis["webpackChunkcontrol"]||[]).push([[879],{3879:(e,t,a)=>{a.r(t),a.d(t,{default:()=>ue});var l=a(9835),s=a(499),n=a(6970),o=a(6647),r=a(8339),i=a(2320),u=a(6397);const c={class:"row justify-start items-center"},m={class:"top-bar-text"},d=(0,l.aZ)({__name:"TopBar",setup(e){const t=(0,o.QT)().t,a=(0,u.n)(),d=[5,10,15,30,60].map((e=>({label:e+"s",value:e}))),p=(0,l.Fl)((()=>{const e=(0,r.yj)();return"/overview"===e.path})),v=(0,l.Fl)((()=>{const e=(0,r.yj)();return e.meta.title?t(e.meta.title.toString()):""})),g=(0,l.Fl)((()=>{const e=(0,r.yj)();return e.meta.icon?e.meta.icon.toString():""})),f=(0,l.Fl)({get:()=>a.getRefreshInterval,set:e=>a.setRefreshInterval(e)});return(e,t)=>{const a=(0,l.up)("q-icon"),o=(0,l.up)("q-space"),r=(0,l.up)("q-toolbar"),u=(0,l.up)("q-header");return(0,l.wg)(),(0,l.j4)(u,{class:"bg-transparent q-pa-sm"},{default:(0,l.w5)((()=>[(0,l.Wm)(r,null,{default:(0,l.w5)((()=>[(0,l._)("div",c,[(0,l.Wm)(a,{color:"primary",name:g.value,size:"md",class:"q-mr-sm"},null,8,["name"]),(0,l._)("h4",m,(0,n.zw)(v.value),1)]),(0,l.Wm)(o),p.value?((0,l.wg)(),(0,l.j4)(i.Z,{key:0,modelValue:f.value,"onUpdate:modelValue":t[0]||(t[0]=e=>f.value=e),options:(0,s.SU)(d),label:e.$t("topBar.selectRefreshRate"),class:"top-bar-refresh-rate-select"},null,8,["modelValue","options","label"])):(0,l.kq)("",!0)])),_:1})])),_:1})}}});var p=a(6602),v=a(1663),g=a(2857),f=a(136),b=a(9984),w=a.n(b);const h=d,y=h;w()(d,"components",{QHeader:p.Z,QToolbar:v.Z,QIcon:g.Z,QSpace:f.Z});var _=a(503),x=a(7747);const q={class:"flex justify-start items-center person-card-box"},M={class:"q-pa-sm",style:{width:"73%"}},Z={class:"q-pb-sm"},I={class:"full-width"},W={class:"flex justify-start items-center"},P=(0,l.aZ)({__name:"ProfileCard",props:{isMiniMenu:{type:Boolean}},setup(e){const t=new _.Z,a=(0,l.Fl)((()=>t.getUsername()));function s(e){return e.charAt(0).toUpperCase()+e.slice(1)}function o(){const e=new x.Z;e.logout()}return(e,t)=>{const r=(0,l.up)("q-avatar"),i=(0,l.up)("q-card"),u=(0,l.up)("q-btn"),c=(0,l.up)("q-tooltip");return!0===e.isMiniMenu?((0,l.wg)(),(0,l.j4)(i,{key:0,flat:"",class:"q-ma-xs q-pa-xs"},{default:(0,l.w5)((()=>[(0,l.Wm)(r,{icon:"person",color:"grey-8","font-size":"36px",size:"xl",class:"text-white"})])),_:1})):((0,l.wg)(),(0,l.j4)(i,{key:1,flat:"",bordered:"",style:{"padding-left":"8px",margin:"4px"}},{default:(0,l.w5)((()=>[(0,l._)("div",q,[(0,l.Wm)(r,{icon:"person",size:"xl",color:"grey-8","font-size":"36px",class:"text-white q-mr-xs"}),(0,l._)("div",M,[(0,l._)("div",Z,(0,n.zw)(s(a.value)),1),(0,l._)("div",I,[(0,l._)("div",W,[(0,l._)("div",null,[(0,l.Wm)(u,{dense:"",color:"primary",icon:"settings",size:"sm",to:"/settings",label:"Settings",class:"q-px-sm q-mr-sm"})]),(0,l._)("div",null,[(0,l.Wm)(u,{dense:"",color:"grey-8",icon:"logout",size:"sm",onClick:t[0]||(t[0]=e=>o())},{default:(0,l.w5)((()=>[(0,l.Wm)(c,{anchor:"bottom middle",class:"bg-grey-8",style:{"font-size":"14px"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(e.$t("profileCard.btnLogout")),1)])),_:1})])),_:1})])])])])])])),_:1}))}}});var k=a(1639),U=a(4458),z=a(1357),Q=a(8879),j=a(6858);const B=(0,k.Z)(P,[["__scopeId","data-v-a7facbd2"]]),S=B;w()(P,"components",{QCard:U.Z,QAvatar:z.Z,QBtn:Q.Z,QTooltip:j.Z});const $={class:"flex justify-evenly q-mt-md",style:{height:"75px"}},C={key:0,class:"text-h6 q-mt-sm"},F={class:"absolute",style:{top:"75px",right:"-15px"}},V=(0,l.aZ)({__name:"SideBar",props:{isVisibleMenu:{type:Boolean}},setup(e){const t=e,a=(0,o.QT)().t,i=(0,s.iH)(t.isVisibleMenu),u=(0,s.iH)(),c=[{menuPosition:4,title:"Gateway",icon:"router",path:"/gateway",disabled:!0,isMenuItem:!0},{menuPosition:5,title:"Backups",icon:"backup",path:"/backups",disabled:!0,isMenuItem:!0},{menuPosition:6,title:"Images",icon:"image",path:"/images",disabled:!0,isMenuItem:!0},{menuPosition:7,title:"CI/CD",icon:"settings_suggest",path:"/ci-cd",disabled:!0,isMenuItem:!1},{menuPosition:8,title:"Security",icon:"security",path:"/security",disabled:!0,isMenuItem:!0},{menuPosition:9,title:"Metrics&Logs",icon:"receipt",path:"/metrics-logs",disabled:!0,isMenuItem:!0},{menuPosition:10,title:"Terminal",icon:"terminal",path:"/terminal",disabled:!0,isMenuItem:!0},{menuPosition:11,title:"IAM",icon:"group",path:"/iam",disabled:!0,isMenuItem:!0},{menuPosition:12,title:"Server Settings",icon:"settings",path:"/server-settings",disabled:!0,isMenuItem:!0}];function m(){const e=(0,r.tv)().getRoutes();let t=[];e.forEach((e=>{const l=e.children||[];l.forEach((e=>{var l,s,n,o,r,i,u,c,m,d;!1!==(null===(l=e.meta)||void 0===l?void 0:l.isMenuItem)&&t.push({menuPosition:null!==(n=null===(s=e.meta)||void 0===s?void 0:s.menuPosition)&&void 0!==n?n:0,title:null!==(r=a(`${null===(o=e.meta)||void 0===o?void 0:o.title}`))&&void 0!==r?r:"",icon:null!==(c=null===(u=null===(i=e.meta)||void 0===i?void 0:i.icon)||void 0===u?void 0:u.toString())&&void 0!==c?c:"",path:e.path,disabled:null!==(d=null===(m=e.meta)||void 0===m?void 0:m.disabled)&&void 0!==d&&d})}))}));let l=t.concat(c);return l.sort(((e,t)=>e.menuPosition<t.menuPosition?-1:e.menuPosition>t.menuPosition?1:0))}return(0,l.YP)((()=>t.isVisibleMenu),(e=>{i.value=e})),(0,l.wF)((()=>{u.value=m()})),(e,t)=>{const o=(0,l.up)("q-icon"),r=(0,l.up)("q-avatar"),c=(0,l.up)("q-btn"),m=(0,l.up)("q-item-section"),d=(0,l.up)("q-item-label"),p=(0,l.up)("q-tooltip"),v=(0,l.up)("q-item"),g=(0,l.up)("q-list"),f=(0,l.up)("q-drawer"),b=(0,l.Q2)("ripple");return(0,l.wg)(),(0,l.j4)(f,{"show-if-above":"",mini:i.value,side:"left",bordered:"",width:220,"mini-width":50},{default:(0,l.w5)((()=>[(0,l.Wm)(g,{padding:""},{default:(0,l.w5)((()=>[(0,l._)("div",$,[(0,l.Wm)(r,null,{default:(0,l.w5)((()=>[(0,l.Wm)(o,{name:"code",size:"xl",color:"primary"})])),_:1}),!1===i.value?((0,l.wg)(),(0,l.iD)("div",C," Speedia Control ")):(0,l.kq)("",!0)]),(0,l.Wm)(S,{isMiniMenu:i.value},null,8,["isMiniMenu"]),(0,l._)("div",F,[(0,l.Wm)(c,{dense:"",round:"",unelevated:"",color:"primary",icon:i.value?"chevron_right":"chevron_left",onClick:t[0]||(t[0]=e=>i.value?i.value=!1:i.value=!0)},null,8,["icon"])]),((0,l.wg)(!0),(0,l.iD)(l.HY,null,(0,l.Ko)(u.value,((e,t)=>(0,l.wy)(((0,l.wg)(),(0,l.j4)(v,{key:t,to:e.path,disable:e.disabled,clickable:""},{default:(0,l.w5)((()=>[e.icon?((0,l.wg)(),(0,l.j4)(m,{key:0,avatar:"",style:{"align-items":"center"}},{default:(0,l.w5)((()=>[(0,l.Wm)(o,{name:e.icon},null,8,["name"])])),_:2},1024)):(0,l.kq)("",!0),(0,l.Wm)(m,null,{default:(0,l.w5)((()=>[(0,l.Wm)(d,null,{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(e.title),1)])),_:2},1024)])),_:2},1024),e.disabled?((0,l.wg)(),(0,l.j4)(p,{key:1,anchor:"center end",class:"bg-primary",style:{"font-size":"14px"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)((0,s.SU)(a)("sideBar.disabledMenu")),1)])),_:1})):(0,l.kq)("",!0)])),_:2},1032,["to","disable"])),[[b]]))),128))])),_:1})])),_:1},8,["mini"])}}});var T=a(906),L=a(3246),H=a(490),R=a(6749),A=a(3115),D=a(1136);const E=V,Y=E;w()(V,"components",{QDrawer:T.Z,QList:L.Z,QAvatar:z.Z,QIcon:g.Z,QBtn:Q.Z,QItem:H.Z,QItemSection:R.Z,QItemLabel:A.Z,QTooltip:j.Z}),w()(V,"directives",{Ripple:D.Z});var G=a(7178);const K={class:"flex justify-end items-center"},J={class:"absolute-full flex flex-center"},N={class:"absolute-full flex flex-center"},O={class:"absolute-full flex flex-center"},X=(0,l.aZ)({__name:"FooterBar",setup(e){const t=new G.Z,a=(0,u.n)(),o=(0,s.iH)(),r=(0,s.iH)(),i=(0,l.Fl)((()=>a.getRefreshInterval));function c(e){return e<50?"green":e<80?"orange":"red"}function m(){o.value&&clearInterval(o.value),o.value=setInterval((()=>{t.getSystemInfo().then((e=>{r.value=e.data.body})).catch((e=>{console.error(e)}))}),1e3*i.value)}return(0,l.wF)((()=>{t.getSystemInfo().then((e=>{r.value=e.data.body})).catch((e=>{console.error(e)})).finally((()=>{m()}))})),(e,t)=>{const a=(0,l.up)("q-tooltip"),s=(0,l.up)("q-icon"),o=(0,l.up)("q-badge"),i=(0,l.up)("q-linear-progress"),u=(0,l.up)("q-footer");return r.value?((0,l.wg)(),(0,l.j4)(u,{key:0,bordered:"",class:"bg-footer q-px-lg"},{default:(0,l.w5)((()=>[(0,l._)("div",K,[(0,l.Wm)(s,{name:"terminal",size:"1.618rem",class:"disabled q-mr-md"},{default:(0,l.w5)((()=>[(0,l.Wm)(a,{class:"bg-primary",style:{"font-size":"14px"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(e.$t("footerBar.disabled")),1)])),_:1})])),_:1}),(0,l.Wm)(s,{name:"memory",size:"sm",class:"q-mr-xs"}),(0,l.Wm)(i,{stripe:"",rounded:"",size:"20px",class:"q-mr-md",value:Math.trunc(r.value.resourceUsage.cpuPercent)/100,color:c(Math.trunc(r.value.resourceUsage.cpuPercent)),label:`${Math.trunc(r.value.resourceUsage.cpuPercent)}%`,style:{width:"100px"}},{default:(0,l.w5)((()=>[(0,l._)("div",J,[(0,l.Wm)(o,{color:"white","text-color":"dark",label:`${Math.trunc(r.value.resourceUsage.cpuPercent)}%`},null,8,["label"])]),(0,l.Wm)(a,{class:"bg-primary",style:{"font-size":"14px"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(e.$t("footerBar.cpuUsage",{cpuUsage:Math.trunc(r.value.resourceUsage.cpuPercent)})),1)])),_:1})])),_:1},8,["value","color","label"]),(0,l.Wm)(s,{name:"memory",size:"sm",class:"q-mr-xs"}),(0,l.Wm)(i,{stripe:"",rounded:"",class:"q-mr-md",size:"20px",value:Math.trunc(r.value.resourceUsage.memoryPercent)/100,color:c(Math.trunc(r.value.resourceUsage.memoryPercent)),label:`${Math.trunc(r.value.resourceUsage.memoryPercent)}%`,style:{width:"100px"}},{default:(0,l.w5)((()=>[(0,l._)("div",N,[(0,l.Wm)(o,{color:"white","text-color":"dark",label:`${Math.trunc(r.value.resourceUsage.memoryPercent)}%`},null,8,["label"])]),(0,l.Wm)(a,{class:"bg-primary",style:{"font-size":"14px"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(e.$t("footerBar.ramUsage",{ramUsage:Math.trunc(r.value.resourceUsage.memoryPercent)})),1)])),_:1})])),_:1},8,["value","color","label"]),(0,l.Wm)(s,{name:"sd_card",size:"sm",class:"q-mr-xs"}),(0,l.Wm)(i,{stripe:"",rounded:"",size:"20px",value:Math.trunc(r.value.resourceUsage.storageInfo[0].usedPercent)/100,color:c(Math.trunc(r.value.resourceUsage.storageInfo[0].usedPercent)),label:`${Math.trunc(r.value.resourceUsage.storageInfo[0].usedPercent)}%`,style:{width:"100px"}},{default:(0,l.w5)((()=>[(0,l._)("div",O,[(0,l.Wm)(o,{color:"white","text-color":"dark",label:`${Math.trunc(r.value.resourceUsage.storageInfo[0].usedPercent)}%`},null,8,["label"])]),(0,l.Wm)(a,{class:"bg-primary",style:{"font-size":"14px"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(e.$t("footerBar.storageInfo",{storageInfo:Math.trunc(r.value.resourceUsage.storageInfo[0].usedPercent)})),1)])),_:1})])),_:1},8,["value","color","label"])])])),_:1})):(0,l.kq)("",!0)}}});var ee=a(1378),te=a(8289),ae=a(990);const le=(0,k.Z)(X,[["__scopeId","data-v-b3d283a0"]]),se=le;w()(X,"components",{QFooter:ee.Z,QIcon:g.Z,QTooltip:j.Z,QLinearProgress:te.Z,QBadge:ae.Z});const ne=(0,l.aZ)({__name:"MainLayout",setup(e){const t=(0,s.iH)(!1);return(e,a)=>{const s=(0,l.up)("router-view"),n=(0,l.up)("q-page-container"),o=(0,l.up)("q-layout");return(0,l.wg)(),(0,l.j4)(o,{view:"lhh Lpr lFf"},{default:(0,l.w5)((()=>[(0,l.Wm)(y,{isVisibleMenu:t.value,"onUpdate:isVisibleMenu":a[0]||(a[0]=e=>t.value=e)},null,8,["isVisibleMenu"]),(0,l.Wm)(Y,{isVisibleMenu:t.value},null,8,["isVisibleMenu"]),(0,l.Wm)(n,null,{default:(0,l.w5)((()=>[(0,l.Wm)(s)])),_:1}),(0,l.Wm)(se)])),_:1})}}});var oe=a(7605),re=a(2133);const ie=ne,ue=ie;w()(ne,"components",{QLayout:oe.Z,QPageContainer:re.Z})}}]);