"use strict";(globalThis["webpackChunkez"]=globalThis["webpackChunkez"]||[]).push([[742],{7089:(e,a,t)=>{t.r(a),t.d(a,{default:()=>Ua});var l=t(9835),o=t(499),s=t(6970),u=t(1957),n=t(4304);const c={class:"q-pb-xs"},i={class:"row justify-between items-center"},r={class:"col-3"},d={class:"col-9 q-pl-md"},p=(0,l.aZ)({__name:"AccountsQuotaInput",props:{quota:{},quotaAvailable:{},title:{},inputLabel:{},step:{},labelValue:{}},emits:["update:quota"],setup(e,{emit:a}){const t=e,u=a,p=(0,o.iH)(0),m=(0,o.iH)(0),g=(0,o.iH)(t.quota),v=(0,o.iH)(!1),w=(0,l.Fl)((()=>p.value<=t.quotaAvailable));return(0,l.YP)(g,(e=>{u("update:quota",e),!0!==v.value&&(p.value=e,m.value++)}),{immediate:!0}),(0,l.YP)(p,(e=>{v.value=!0,g.value=e,setTimeout((()=>{v.value=!1}),10)}),{immediate:!0}),(e,a)=>{var o,u,v;const b=(0,l.up)("q-slider");return(0,l.wg)(),(0,l.iD)(l.HY,null,[(0,l._)("div",c,(0,s.zw)(e.title),1),(0,l._)("div",i,[(0,l._)("div",r,[((0,l.wg)(),(0,l.j4)(n.Z,{label:e.inputLabel,value:p.value,"onUpdate:value":a[0]||(a[0]=e=>p.value=e),key:m.value,rules:[()=>w.value||e.$t("accountsQuotaInput.invalidValue")],suffix:null!==(o=t.labelValue)&&void 0!==o?o:"GB",maxLength:"5"},null,8,["label","value","rules","suffix"]))]),(0,l._)("div",d,[(0,l.Wm)(b,{modelValue:g.value,"onUpdate:modelValue":a[1]||(a[1]=e=>g.value=e),disable:0===e.quotaAvailable,min:0,max:e.quotaAvailable,step:null!==(u=t.step)&&void 0!==u?u:1,label:"","label-value":`${g.value} ${null!==(v=t.labelValue)&&void 0!==v?v:"GB"}`,"switch-label-side":"","label-always":"",markers:""},null,8,["modelValue","disable","max","step","label-value"])])])],64)}}});var m=t(8423),g=t(9984),v=t.n(g);const w=p,b=w;v()(p,"components",{QSlider:m.Z});var y=t(9476);const f={class:"q-pb-xs"},Q={class:"row justify-between items-center"},U={class:"col-3"},_={class:"col-9 q-pl-md"},q=(0,l.aZ)({__name:"AccountsInodesQuota",props:{inodesQuota:{},inodesAvailable:{}},emits:["update:inodesQuota"],setup(e,{emit:a}){const t=e,u=a,c=(0,o.iH)(t.inodesQuota),i=(0,o.iH)(0),r=(0,o.iH)(0),d=(0,o.iH)(!1),p=(0,l.Fl)((()=>i.value<=t.inodesAvailable));return(0,l.YP)(c,(e=>{u("update:inodesQuota",e),!0!==d.value&&(i.value=e,r.value++)}),{immediate:!0}),(0,l.YP)(i,(e=>{d.value=!0,i.value;{const a=e.toString();i.value=parseInt(a)}c.value=i.value,setTimeout((()=>{d.value=!1}),10)}),{immediate:!0}),(e,a)=>{const t=(0,l.up)("q-slider");return(0,l.wg)(),(0,l.iD)(l.HY,null,[(0,l._)("div",f,(0,s.zw)(e.$t("accountsUpdateQuotaDialog.labelInodesQuota")),1),(0,l._)("div",Q,[(0,l._)("div",U,[((0,l.wg)(),(0,l.j4)(n.Z,{label:e.$t("accountsUpdateQuotaDialog.labelInodesQuota"),value:i.value,"onUpdate:value":a[0]||(a[0]=e=>i.value=e),key:r.value,rules:[()=>p.value||e.$t("accountsUpdateQuotaDialog.invalidValue")]},null,8,["label","value","rules"]))]),(0,l._)("div",_,[(0,l.Wm)(t,{modelValue:c.value,"onUpdate:modelValue":a[1]||(a[1]=e=>c.value=e),disable:0===e.inodesAvailable,min:0,max:e.inodesAvailable,label:"","label-value":(0,o.SU)(y.Z)(c.value,!0,2),"switch-label-side":"","label-always":"",markers:""},null,8,["modelValue","disable","max","label-value"])])])],64)}}}),D=q,h=D;v()(q,"components",{QSlider:m.Z});var A=t(6397);const Z={class:"q-pb-lg"},x={class:"q-pb-xs"},W={class:"q-py-sm q-pr-sm"},C={class:"q-py-md q-pr-sm"},k={class:"q-pb-md q-pr-sm"},P={class:"q-py-sm q-pr-sm"},S=["innerHTML"],F=128,B=128,H=4096,$=1e8,T=100,K=(0,l.aZ)({__name:"AccountsQuotasForm",props:{cpuQuota:{},memoryQuota:{},storageQuota:{},inodesQuota:{},storagePerformanceQuota:{}},emits:["update:cpuQuota","update:memoryQuota","update:storageQuota","update:inodesQuota","update:storagePerformanceQuota"],setup(e,{emit:a}){const t=e,u=a,n=(0,A.n)(),c=(0,o.iH)(t.cpuQuota),i=(0,o.iH)(t.memoryQuota),r=(0,o.iH)(t.storageQuota),d=(0,o.iH)(t.inodesQuota),p=(0,o.iH)(t.storagePerformanceQuota),m=(0,l.Fl)((()=>n.getSystemInfo)),g=(0,l.Fl)((()=>{const e={freeBytes:0,totalInodes:0};return m.value.resourceUsage.storageInfo.forEach((a=>{"xfs"===a.fileSystem&&(e.freeBytes=H,e.totalInodes=$)})),e})),v=(0,l.Fl)((()=>({cpuCores:F,memoryBytes:B,storageBytes:g.value.freeBytes,storageInodes:g.value.totalInodes,storagePerformanceUnits:T})));return(0,l.YP)((()=>t.cpuQuota),(e=>{c.value=e})),(0,l.YP)((()=>t.memoryQuota),(e=>{i.value=e})),(0,l.YP)((()=>t.storageQuota),(e=>{r.value=e})),(0,l.YP)((()=>t.inodesQuota),(e=>{d.value=e})),(0,l.YP)((()=>t.storagePerformanceQuota),(e=>{p.value=e})),(0,l.YP)(c,(e=>{u("update:cpuQuota",e)})),(0,l.YP)(i,(e=>{u("update:memoryQuota",e)})),(0,l.YP)(r,(e=>{u("update:storageQuota",e)})),(0,l.YP)(d,(e=>{u("update:inodesQuota",e)})),(0,l.YP)(p,(e=>{u("update:storagePerformanceQuota",e)})),(e,a)=>{const t=(0,l.up)("q-slider"),o=(0,l.up)("q-banner");return(0,l.wg)(),(0,l.iD)(l.HY,null,[(0,l._)("div",Z,[(0,l._)("div",x,(0,s.zw)(e.$t("accountsQuotasForm.cpuQuotaTitle")),1),(0,l.Wm)(t,{modelValue:c.value,"onUpdate:modelValue":a[0]||(a[0]=e=>c.value=e),disable:0===v.value.cpuCores,min:0,step:4,max:F,label:"","label-value":c.value+" Cores","switch-label-side":"","label-always":"",markers:"",class:"q-px-xs"},null,8,["modelValue","disable","label-value"])]),(0,l._)("div",W,[(0,l.Wm)(b,{quota:r.value,"onUpdate:quota":a[1]||(a[1]=e=>r.value=e),quotaAvailable:H,title:e.$t("accountsQuotasForm.storageQuotaTitle"),inputLabel:e.$t("accountsQuotasForm.storageQuotaInputLabel")},null,8,["quota","title","inputLabel"])]),(0,l._)("div",C,[(0,l.Wm)(b,{quota:i.value,"onUpdate:quota":a[2]||(a[2]=e=>i.value=e),quotaAvailable:B,title:e.$t("accountsQuotasForm.memoryQuotaTitle"),step:8,inputLabel:e.$t("accountsQuotasForm.memoryQuotaInputLabel")},null,8,["quota","title","inputLabel"])]),(0,l._)("div",k,[(0,l.Wm)(h,{inodesQuota:d.value,"onUpdate:inodesQuota":a[3]||(a[3]=e=>d.value=e),inodesAvailable:$},null,8,["inodesQuota"])]),(0,l._)("div",P,[(0,l.Wm)(b,{quota:p.value,"onUpdate:quota":a[4]||(a[4]=e=>p.value=e),quotaAvailable:T,title:e.$t("accountsQuotasForm.storagePerformanceQuotaTitle"),inputLabel:e.$t("accountsQuotasForm.storagePerformanceQuotaInputLabel"),labelValue:"Units"},null,8,["quota","title","inputLabel"])]),(0,l.Wm)(o,{class:"dialog-info"},{default:(0,l.w5)((()=>[(0,l._)("div",{innerHTML:e.$t("accountsQuotasForm.aboutStoragePerformanceUnit")},null,8,S)])),_:1})],64)}}});var z=t(7128);const I=K,V=I;v()(K,"components",{QSlider:m.Z,QBanner:z.Z});var L=t(5786),j=t(2626),Y=t(7713),G=t(3064),M=t(5971),X=t(8900),N=t(3055),E=t(6647);const R={class:"flex justify-between items-center q-mb-sm"},J={class:"title-dialog"},O=(0,l.aZ)({__name:"AccountsCreateUserDialog",props:{showCreateUserDialog:{type:Boolean}},emits:["update:showCreateUserDialog"],setup(e,{emit:a}){const t=e,c=a,i=(0,E.QT)().t,r=(0,N.o)(),d=(0,o.iH)(t.showCreateUserDialog),p=(0,o.iH)(""),m=(0,o.iH)(""),g=(0,o.iH)(0),v=(0,o.iH)(0),w=(0,o.iH)(0),b=(0,o.iH)(0),y=(0,o.iH)(0),f=(0,o.iH)(!1),Q=(0,l.Fl)((()=>p.value.length>=3)),U=(0,l.Fl)({get:()=>r.getKeyAccountsTable,set:e=>{r.setKeyAccountsTable(e)}});function _(){d.value=!1}(0,l.YP)((()=>t.showCreateUserDialog),(e=>{d.value=e})),(0,l.YP)(d,(e=>{c("update:showCreateUserDialog",e)}));const q=Math.pow(1024,3);function D(){const e=new M.Z;e.createAccount({username:p.value,password:m.value,quota:{cpuCores:g.value,memoryBytes:v.value*q,storageBytes:w.value*q,storageInodes:b.value,storagePerformanceUnits:y.value}}).then((()=>{(0,X.LX)(i("accountsCreateUserDialog.createdSuccessfully")),U.value++,_()})).catch((e=>{(0,X.s9)(e.response.data,i("accountsCreateUserDialog.errorCreatingAccount"))}))}return(e,a)=>{const t=(0,l.up)("q-btn"),o=(0,l.up)("q-card-section"),c=(0,l.up)("q-card-actions"),i=(0,l.up)("q-card"),r=(0,l.up)("q-dialog");return(0,l.wg)(),(0,l.j4)(r,{modelValue:d.value,"onUpdate:modelValue":a[12]||(a[12]=e=>d.value=e),persistent:""},{default:(0,l.w5)((()=>[(0,l.Wm)(i,{class:"dialog-card-bg",style:{width:"700px","max-width":"80vw"}},{default:(0,l.w5)((()=>[(0,l._)("div",R,[(0,l._)("div",J,(0,s.zw)(e.$t("accountsCreateUserDialog.title")),1),(0,l.Wm)(t,{flat:"",round:"",dense:"",icon:"close",onClick:a[0]||(a[0]=e=>_())})]),(0,l.Wm)(o,{class:"q-px-xs"},{default:(0,l.w5)((()=>[(0,l.Wm)(n.Z,{label:e.$t("accountsCreateUserDialog.username"),icon:"person",value:p.value,"onUpdate:value":a[1]||(a[1]=e=>p.value=e),rules:[()=>""!==p.value||e.$t("accountsCreateUserDialog.usernameRequired"),()=>p.value.length>=3||e.$t("accountsCreateUserDialog.usernameMinLength")]},null,8,["label","value","rules"])])),_:1}),(0,l.Wm)(o,{class:"q-px-xs"},{default:(0,l.w5)((()=>[(0,l.Wm)(j.Z,{password:m.value,"onUpdate:password":a[2]||(a[2]=e=>m.value=e)},null,8,["password"]),(0,l.Wm)(Y.Z,{password:m.value,"onUpdate:password":a[3]||(a[3]=e=>m.value=e),class:"float-right"},null,8,["password"])])),_:1}),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.wy)((0,l.Wm)(G.Z,{password:m.value,minPasswordLength:6,"onUpdate:isValidPassword":a[4]||(a[4]=e=>f.value=e)},null,8,["password"]),[[u.F8,m.value.length>0]])])),_:1}),(0,l.Wm)(o,{class:"q-px-md"},{default:(0,l.w5)((()=>[(0,l.Wm)(V,{cpuQuota:g.value,"onUpdate:cpuQuota":a[5]||(a[5]=e=>g.value=e),memoryQuota:v.value,"onUpdate:memoryQuota":a[6]||(a[6]=e=>v.value=e),storageQuota:w.value,"onUpdate:storageQuota":a[7]||(a[7]=e=>w.value=e),inodesQuota:b.value,"onUpdate:inodesQuota":a[8]||(a[8]=e=>b.value=e),storagePerformanceQuota:y.value,"onUpdate:storagePerformanceQuota":a[9]||(a[9]=e=>y.value=e)},null,8,["cpuQuota","memoryQuota","storageQuota","inodesQuota","storagePerformanceQuota"])])),_:1}),(0,l.Wm)(c,{align:"between",class:"q-pt-md q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(L.Z,{label:e.$t("accountsCreateUserDialog.cancelBtn"),class:"q-mr-sm q-px-xl",color:"grey-8",onClick:a[10]||(a[10]=e=>_())},null,8,["label"]),(0,l.Wm)(L.Z,{color:"primary",icon:"person_add",label:e.$t("accountsCreateUserDialog.createBtn"),disable:!Q.value||!f.value,onClick:a[11]||(a[11]=e=>D())},null,8,["label","disable"])])),_:1})])),_:1})])),_:1},8,["modelValue"])}}});var ee=t(2074),ae=t(4458),te=t(8879),le=t(3190),oe=t(1821);const se=O,ue=se;v()(O,"components",{QDialog:ee.Z,QCard:ae.Z,QBtn:te.Z,QCardSection:le.Z,QCardActions:oe.Z});const ne={class:"flex justify-end items-center"},ce={class:"td-action-bg-15"},ie={class:"td-action-bg-20"},re={class:"td-action-bg-25"},de=(0,l.aZ)({__name:"AccountTableActions",props:{account:{}},setup(e){const a=e,t=(0,N.o)(),o=(0,l.Fl)({get:()=>t.getKeyUpdatePasswordDialog,set:e=>{t.setKeyUpdatePasswordDialog(e)}}),u=(0,l.Fl)({get:()=>t.getKeyUpdateApiKeyDialog,set:e=>{t.setKeyUpdateApiKeyDialog(e)}}),n=(0,l.Fl)({get:()=>t.getKeyUpdateQuotaDialog,set:e=>{t.setKeyUpdateQuotaDialog(e)}});function c(){o.value++,t.setAccountSelected(a.account),t.setShowUpdatePasswordDialog(!0)}function i(){t.setAccountSelected(a.account),t.setShowDeleteUserDialog(!0)}function r(){u.value++,t.setAccountSelected(a.account),t.setShowUpdateApiKeyDialog(!0)}function d(){n.value++,t.setAccountSelected(a.account),t.setShowUpdateQuotaDialog(!0)}return(e,a)=>{const t=(0,l.up)("q-tooltip"),o=(0,l.up)("q-btn");return(0,l.wg)(),(0,l.iD)("div",ne,[(0,l._)("div",ce,[(0,l.Wm)(o,{size:"md",icon:"lock",flat:"",onClick:a[0]||(a[0]=e=>c())},{default:(0,l.w5)((()=>[(0,l.Wm)(t,{class:"bg-primary text-white",style:{"font-size":"14px"},offset:[10,10]},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(e.$t("accountsTableActions.editPasswordBtn")),1)])),_:1})])),_:1})]),(0,l._)("div",ie,[(0,l.Wm)(o,{size:"md",icon:"key",flat:"",onClick:a[1]||(a[1]=e=>r())},{default:(0,l.w5)((()=>[(0,l.Wm)(t,{class:"bg-primary text-white",style:{"font-size":"14px"},offset:[10,10]},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(e.$t("accountsTableActions.editApiKeyBtn")),1)])),_:1})])),_:1})]),(0,l._)("div",re,[(0,l.Wm)(o,{size:"md",icon:"settings",flat:"",onClick:a[2]||(a[2]=e=>d())},{default:(0,l.w5)((()=>[(0,l.Wm)(t,{class:"bg-primary text-white",style:{"font-size":"14px"},offset:[10,10]},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(e.$t("accountsTableActions.editQuotaBtn")),1)])),_:1})])),_:1})]),(0,l._)("div",null,[(0,l.Wm)(o,{color:"negative",icon:"delete",size:"md",onClick:a[3]||(a[3]=e=>i())},{default:(0,l.w5)((()=>[(0,l.Wm)(t,{class:"bg-negative text-white",style:{"font-size":"14px"},offset:[10,10]},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(e.$t("accountsTableActions.deleteBtn")),1)])),_:1})])),_:1})])])}}});var pe=t(6858);const me=de,ge=me;v()(de,"components",{QBtn:te.Z,QTooltip:pe.Z});var ve=t(8428),we=t(7250);const be=e=>((0,l.dD)("data-v-35824e3f"),e=e(),(0,l.Cn)(),e),ye=be((()=>(0,l._)("th",{class:"text-center account-table-th"},"CPU",-1))),fe=be((()=>(0,l._)("th",{class:"text-center account-table-th"},"Memory",-1))),Qe=be((()=>(0,l._)("th",{class:"text-center account-table-th"},"Disk",-1))),Ue=be((()=>(0,l._)("th",{class:"text-center account-table-th"},"Inodes",-1))),_e=be((()=>(0,l._)("th",{class:"text-center account-table-th"},"Storage Performance Units",-1))),qe=be((()=>(0,l._)("th",{class:"text-center account-table-th"},"UserId / GroupId",-1))),De={style:{width:"150px"}},he={style:{width:"150px"}},Ae={style:{width:"150px"}},Ze={style:{width:"150px"}},xe={class:"text-center"},We={class:"text-center"},Ce={class:"row justify-center q-mt-md"},ke=(0,l.aZ)({__name:"AccountsTable",props:{accountsList:{}},setup(e){const a=e,t=(0,o.iH)(!1),u=(0,o.iH)(0),n=(0,o.iH)(""),c=(0,o.iH)([{name:"username",label:"Username",align:"left",field:"username",classes:"td-main-table",headerClasses:"bg-primary text-white",sortable:!0}]),i=(0,o.iH)({sortBy:"desc",descending:!1,page:1,rowsPerPage:10}),r=(0,l.Fl)((()=>Math.ceil(a.accountsList.length/i.value.rowsPerPage))),d=e=>parseFloat((0,we.Z)(e,!1,2,"GiB"));function p(){u.value++,t.value=!0}return(e,a)=>{const o=(0,l.up)("q-icon"),m=(0,l.up)("q-input"),g=(0,l.up)("q-space"),v=(0,l.up)("q-th"),w=(0,l.up)("q-tr"),b=(0,l.up)("q-td"),y=(0,l.up)("q-table"),f=(0,l.up)("q-pagination");return(0,l.wg)(),(0,l.iD)(l.HY,null,[e.accountsList.length>0?((0,l.wg)(),(0,l.j4)(ue,{showCreateUserDialog:t.value,"onUpdate:showCreateUserDialog":a[0]||(a[0]=e=>t.value=e),key:u.value},null,8,["showCreateUserDialog"])):(0,l.kq)("",!0),(0,l.Wm)(y,{rows:e.accountsList,columns:c.value,filter:n.value,pagination:i.value,"onUpdate:pagination":a[3]||(a[3]=e=>i.value=e),"no-data-label":e.$t("accountsTable.noDataLabel"),"row-key":"key",color:"primary",flat:"",bordered:"","hide-pagination":""},{top:(0,l.w5)((()=>[(0,l.Wm)(m,{borderless:"",debounce:"300",color:"primary",modelValue:n.value,"onUpdate:modelValue":a[1]||(a[1]=e=>n.value=e),label:e.$t("accountsTable.searchInput")},{prepend:(0,l.w5)((()=>[(0,l.Wm)(o,{name:"search"})])),_:1},8,["modelValue","label"]),(0,l.Wm)(g),(0,l.Wm)(L.Z,{label:e.$t("accountsTable.createAccountBtn"),color:"primary",icon:"person_add",onClick:a[2]||(a[2]=e=>p())},null,8,["label"])])),header:(0,l.w5)((e=>[(0,l.Wm)(w,{props:e},{default:(0,l.w5)((()=>[((0,l.wg)(!0),(0,l.iD)(l.HY,null,(0,l.Ko)(e.cols,(a=>((0,l.wg)(),(0,l.j4)(v,{key:a.name,props:e,style:{"font-weight":"bold","font-size":"14px"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(a.label),1)])),_:2},1032,["props"])))),128)),ye,fe,Qe,Ue,_e,qe,(0,l.Wm)(v)])),_:2},1032,["props"])])),body:(0,l.w5)((e=>[(0,l.Wm)(w,{props:e},{default:(0,l.w5)((()=>[((0,l.wg)(!0),(0,l.iD)(l.HY,null,(0,l.Ko)(e.cols,(a=>((0,l.wg)(),(0,l.j4)(b,{key:a.name,props:e},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(a.value),1)])),_:2},1032,["props"])))),128)),(0,l._)("td",De,[(0,l.Wm)(ve.Z,{total:e.row.quota.cpuCores,progress:e.row.quotaUsage.cpuCores,suffix:1===e.row.quota.cpuCores?"Core":"Cores"},null,8,["total","progress","suffix"])]),(0,l._)("td",he,[(0,l.Wm)(ve.Z,{total:d(e.row.quota.memoryBytes),progress:d(e.row.quotaUsage.memoryBytes),suffix:"GB"},null,8,["total","progress"])]),(0,l._)("td",Ae,[(0,l.Wm)(ve.Z,{total:d(e.row.quota.storageBytes),progress:d(e.row.quotaUsage.storageBytes),suffix:"GB"},null,8,["total","progress"])]),(0,l._)("td",Ze,[(0,l.Wm)(ve.Z,{total:d(e.row.quota.storageBytes),progress:d(e.row.quotaUsage.storageBytes),suffix:"GB"},null,8,["total","progress"])]),(0,l._)("td",xe,(0,s.zw)(e.row.quota.storagePerformanceUnits),1),(0,l._)("td",We,(0,s.zw)(e.row.id)+" / "+(0,s.zw)(e.row.groupId),1),(0,l.Wm)(b,{class:"text-right"},{default:(0,l.w5)((()=>[(0,l.Wm)(ge,{account:e.row},null,8,["account"])])),_:2},1024)])),_:2},1032,["props"])])),_:1},8,["rows","columns","filter","pagination","no-data-label"]),(0,l._)("div",Ce,[(0,l.Wm)(f,{modelValue:i.value.page,"onUpdate:modelValue":a[4]||(a[4]=e=>i.value.page=e),color:"primary",max:r.value,size:"md"},null,8,["modelValue","max"])])],64)}}});var Pe=t(1639),Se=t(422),Fe=t(6611),Be=t(2857),He=t(136),$e=t(1233),Te=t(1682),Ke=t(7220),ze=t(996);const Ie=(0,Pe.Z)(ke,[["__scopeId","data-v-35824e3f"]]),Ve=Ie;v()(ke,"components",{QTable:Se.Z,QInput:Fe.Z,QIcon:Be.Z,QSpace:He.Z,QTr:$e.Z,QTh:Te.Z,QTd:Ke.Z,QPagination:ze.Z});var Le=t(5273);const je={class:"flex justify-between items-center"},Ye={class:"title-dialog"},Ge=(0,l.aZ)({__name:"AccountsUpdatePasswordDialog",setup(e){const a=(0,E.QT)().t,t=(0,N.o)(),n=(0,l.Fl)({get:()=>t.getShowUpdatePasswordDialog,set:e=>{t.setShowUpdatePasswordDialog(e)}}),c=(0,l.Fl)((()=>t.getAccountSelected.id)),i=(0,l.Fl)((()=>t.getAccountSelected.username)),r=(0,o.iH)(""),d=(0,o.iH)(!1),p=(0,l.Fl)({get:()=>t.keyAccountsTable,set:e=>{t.keyAccountsTable=e}});function m(){n.value=!1}function g(){(0,Le.Q)();const e=new M.Z;e.updateAccount({accountId:c.value,password:r.value}).then((()=>{(0,X.LX)(a("accountsUpdatePasswordDialog.updatedSuccessfully")),p.value++,m()})).catch((e=>{(0,X.s9)(e.response.data,a("accountsUpdatePasswordDialog.errorUpdatingAccount"))})).finally((()=>{(0,Le.Z)()}))}return(e,a)=>{const t=(0,l.up)("q-btn"),o=(0,l.up)("q-card-section"),c=(0,l.up)("q-card-actions"),p=(0,l.up)("q-card"),v=(0,l.up)("q-dialog");return(0,l.wg)(),(0,l.j4)(v,{modelValue:n.value,"onUpdate:modelValue":a[6]||(a[6]=e=>n.value=e),persistent:""},{default:(0,l.w5)((()=>[(0,l.Wm)(p,{style:{width:"700px","max-width":"80vw"},class:"dialog-card-bg"},{default:(0,l.w5)((()=>[(0,l._)("div",je,[(0,l._)("div",Ye,(0,s.zw)(e.$t("accountsUpdatePasswordDialog.title",{username:i.value})),1),(0,l.Wm)(t,{flat:"",round:"",dense:"",icon:"close",onClick:a[0]||(a[0]=e=>m())})]),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(j.Z,{password:r.value,"onUpdate:password":a[1]||(a[1]=e=>r.value=e),minLength:6},null,8,["password"]),(0,l.Wm)(Y.Z,{password:r.value,"onUpdate:password":a[2]||(a[2]=e=>r.value=e),class:"float-right"},null,8,["password"])])),_:1}),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.wy)((0,l.Wm)(G.Z,{minPasswordLength:6,password:r.value,"onUpdate:isValidPassword":a[3]||(a[3]=e=>d.value=e)},null,8,["password"]),[[u.F8,r.value.length>0]])])),_:1}),(0,l.Wm)(c,{align:"between",class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(L.Z,{label:e.$t("accountsUpdatePasswordDialog.cancelBtn"),color:"grey-8",onClick:a[4]||(a[4]=e=>m())},null,8,["label"]),(0,l.Wm)(L.Z,{color:"primary",icon:"key",label:e.$t("accountsUpdatePasswordDialog.updateBtn"),disable:!1===d.value,onClick:a[5]||(a[5]=e=>g())},null,8,["label","disable"])])),_:1})])),_:1})])),_:1},8,["modelValue"])}}}),Me=Ge,Xe=Me;v()(Ge,"components",{QDialog:ee.Z,QCard:ae.Z,QBtn:te.Z,QCardSection:le.Z,QCardActions:oe.Z});const Ne={class:"flex justify-between items-center"},Ee={class:"title-dialog"},Re={key:0},Je={class:"box-api-key q-mb-sm"},Oe={class:"flex justify-between"},ea={class:"text-primary",style:{"font-size":"16px"}},aa={style:{"font-size":"14px"}},ta={key:1,style:{height:"50px"}},la=(0,l.aZ)({__name:"AccountsUpdateApiKeyDialog",setup(e){const a=(0,E.QT)().t,t=(0,N.o)(),u=(0,o.iH)(""),n=(0,l.Fl)({get:()=>t.getShowUpdateApiKeyDialog,set:e=>{t.setShowUpdateApiKeyDialog(e)}}),c=(0,l.Fl)((()=>t.getAccountSelected.id)),i=(0,l.Fl)((()=>t.getAccountSelected.username));function r(e){navigator.clipboard.writeText(e),(0,X.c0)({msg:`${a("accountsUpdateApiKeyDialog.copiedToClipboard")}`,position:"bottom",type:"primary",html:!0})}function d(){n.value=!1}function p(){(0,Le.Q)();const e=new M.Z;e.updateAccount({accountId:c.value,shouldUpdateApiKey:!0}).then((e=>{u.value=e.data.body,(0,X.LX)(a("accountsUpdateApiKeyDialog.updatedSuccessfullyWithApiKey"))})).catch((e=>{(0,X.s9)(e.response.data,a("accountsUpdateApiKeyDialog.errorUpdatingAccount"))})).finally((()=>{(0,Le.Z)()}))}return(e,a)=>{const t=(0,l.up)("q-btn"),o=(0,l.up)("q-icon"),c=(0,l.up)("q-card-section"),m=(0,l.up)("q-card-actions"),g=(0,l.up)("q-card"),v=(0,l.up)("q-dialog");return(0,l.wg)(),(0,l.j4)(v,{modelValue:n.value,"onUpdate:modelValue":a[4]||(a[4]=e=>n.value=e),persistent:""},{default:(0,l.w5)((()=>[(0,l.Wm)(g,{style:{width:"700px","max-width":"80vw"},class:"dialog-card-bg"},{default:(0,l.w5)((()=>[(0,l._)("div",Ne,[(0,l._)("div",Ee,(0,s.zw)(e.$t("accountsUpdateApiKeyDialog.title",{username:i.value})),1),(0,l.Wm)(t,{flat:"",round:"",dense:"",icon:"close",onClick:a[0]||(a[0]=e=>d())})]),(0,l.Wm)(c,{class:"q-px-none"},{default:(0,l.w5)((()=>[u.value?((0,l.wg)(),(0,l.iD)("div",Re,[(0,l._)("div",Je,[(0,l._)("div",Oe,[(0,l._)("p",ea,(0,s.zw)(e.$t("accountsUpdateApiKeyDialog.generatedApiKey")),1),(0,l.Wm)(o,{name:"content_copy",size:"20px",class:"cursor-pointer icon-copy-api-key",onClick:a[1]||(a[1]=e=>r(u.value))})]),(0,l.Uk)(" "+(0,s.zw)(u.value),1)]),(0,l._)("small",aa,[(0,l.Wm)(o,{color:"amber",size:"sm",name:"warning"}),(0,l.Uk)(" "+(0,s.zw)(e.$t("accountsUpdateApiKeyDialog.saveNewApiKey")),1)])])):((0,l.wg)(),(0,l.iD)("div",ta))])),_:1}),(0,l.Wm)(m,{align:"between",class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(L.Z,{label:e.$t("accountsUpdateApiKeyDialog.closeDialogBtn"),color:"grey-8",onClick:a[2]||(a[2]=e=>d())},null,8,["label"]),(0,l.Wm)(L.Z,{icon:"key",label:e.$t("accountsUpdateApiKeyDialog.generateNewApiKeyBtn"),color:"primary",onClick:a[3]||(a[3]=e=>p())},null,8,["label"])])),_:1})])),_:1})])),_:1},8,["modelValue"])}}}),oa=la,sa=oa;v()(la,"components",{QDialog:ee.Z,QCard:ae.Z,QBtn:te.Z,QCardSection:le.Z,QIcon:Be.Z,QCardActions:oe.Z});var ua=t(9906),na=t(9302);const ca=(0,l.aZ)({__name:"AccountsDeleteUserDialog",setup(e){const a=(0,na.Z)(),t=(0,E.QT)().t,o=(0,N.o)(),s=(0,l.Fl)((()=>a.dark.isActive?"/_/icons/bomb_dark.svg":"/_/icons/bomb_light.svg")),u=(0,l.Fl)({get:()=>o.getShowDeleteUserDialog,set:e=>{o.setShowDeleteUserDialog(e)}}),n=(0,l.Fl)((()=>o.getAccountSelected.username)),c=(0,l.Fl)((()=>o.getAccountSelected.id)),i=(0,l.Fl)({get:()=>o.keyAccountsTable,set:e=>{o.keyAccountsTable=e}});function r(){u.value=!1}function d(){(0,Le.Q)();const e=new M.Z;e.deleteAccount(c.value).then((()=>{(0,X.LX)(t("accountsDeleteUserDialog.deletedSuccessfully",{username:n.value})),i.value++,r()})).catch((e=>{(0,X.s9)(e.response.data,t("accountsDeleteUserDialog.errorDeletingAccount",{username:n.value}))})).finally((()=>{(0,Le.Z)()}))}return(e,a)=>((0,l.wg)(),(0,l.j4)(ua.Z,{showDeleteDialog:u.value,"onUpdate:showDeleteDialog":a[2]||(a[2]=e=>u.value=e),titleDialog:e.$t("accountsDeleteUserDialog.title",{username:n.value}),imagePath:s.value,messageToDelete:e.$t("accountsDeleteUserDialog.messageDeleteAccount",{username:n.value}),warningToDelete:e.$t("accountsDeleteUserDialog.warningDeleteAccount")},{"card-actions":(0,l.w5)((()=>[(0,l.Wm)(L.Z,{label:e.$t("accountsDeleteUserDialog.cancelBtn"),color:"grey-8",onClick:a[0]||(a[0]=e=>r())},null,8,["label"]),(0,l.Wm)(L.Z,{color:"negative",label:e.$t("accountsDeleteUserDialog.deleteBtn"),onClick:a[1]||(a[1]=e=>d())},null,8,["label"])])),_:1},8,["showDeleteDialog","titleDialog","imagePath","messageToDelete","warningToDelete"]))}}),ia=ca,ra=ia,da={class:"flex justify-between items-center"},pa={class:"title-dialog"},ma=(0,l.aZ)({__name:"AccountsUpdateQuotaDialog",setup(e){const a=(0,E.QT)().t,t=(0,N.o)(),u=(0,o.iH)(0),c=(0,o.iH)(0),i=(0,o.iH)(0),r=(0,o.iH)(0),d=(0,o.iH)(0),p=(0,l.Fl)((()=>t.getAccountSelected.id)),m=(0,l.Fl)((()=>t.getAccountSelected)),g=(0,l.Fl)({get:()=>t.getKeyAccountsTable,set:e=>{t.setKeyAccountsTable(e)}}),v=(0,l.Fl)({get:()=>t.getShowUpdateQuotaDialog,set:e=>{t.setShowUpdateQuotaDialog(e)}}),w=(0,l.Fl)((()=>t.getAccountSelected.username)),b=Math.pow(1024,3);function y(){v.value=!1}function f(){const e=new M.Z;e.updateAccount({accountId:p.value,quota:{cpuCores:u.value,memoryBytes:c.value*b,storageBytes:i.value*b,storageInodes:r.value,storagePerformanceUnits:d.value}}).then((()=>{(0,X.LX)(a("accountsUpdateQuotaDialog.updatedSuccessfully")),g.value++,y()})).catch((e=>{(0,X.s9)(e.response.data,a("accountsUpdateQuotaDialog.errorUpdatingQuota"))}))}return(0,l.YP)(v,(e=>{!1!==e&&(u.value=m.value.quota.cpuCores,c.value=m.value.quota.memoryBytes/b,i.value=m.value.quota.storageBytes/b,r.value=m.value.quota.storageInodes,d.value=m.value.quota.storagePerformanceUnits)}),{immediate:!0}),(e,a)=>{const t=(0,l.up)("q-btn"),o=(0,l.up)("q-card-section"),p=(0,l.up)("q-card-actions"),m=(0,l.up)("q-card"),g=(0,l.up)("q-dialog");return(0,l.wg)(),(0,l.j4)(g,{modelValue:v.value,"onUpdate:modelValue":a[9]||(a[9]=e=>v.value=e),persistent:""},{default:(0,l.w5)((()=>[(0,l.Wm)(m,{class:"dialog-card-bg",style:{width:"700px","max-width":"80vw","overflow-x":"hidden"}},{default:(0,l.w5)((()=>[(0,l._)("div",da,[(0,l._)("div",pa,(0,s.zw)(e.$t("accountsUpdateQuotaDialog.title",{username:w.value})),1),(0,l.Wm)(t,{flat:"",round:"",dense:"",icon:"close",onClick:a[0]||(a[0]=e=>y())})]),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(n.Z,{label:e.$t("accountsUpdateQuotaDialog.username"),disable:"",icon:"person",value:w.value,"onUpdate:value":a[1]||(a[1]=e=>w.value=e)},null,8,["label","value"])])),_:1}),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(V,{cpuQuota:u.value,"onUpdate:cpuQuota":a[2]||(a[2]=e=>u.value=e),memoryQuota:c.value,"onUpdate:memoryQuota":a[3]||(a[3]=e=>c.value=e),storageQuota:i.value,"onUpdate:storageQuota":a[4]||(a[4]=e=>i.value=e),inodesQuota:r.value,"onUpdate:inodesQuota":a[5]||(a[5]=e=>r.value=e),storagePerformanceQuota:d.value,"onUpdate:storagePerformanceQuota":a[6]||(a[6]=e=>d.value=e)},null,8,["cpuQuota","memoryQuota","storageQuota","inodesQuota","storagePerformanceQuota"])])),_:1}),(0,l.Wm)(p,{align:"between",class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(L.Z,{label:e.$t("accountsUpdateQuotaDialog.cancelBtn"),class:"q-mr-sm q-px-xl",color:"grey-8",onClick:a[7]||(a[7]=e=>y())},null,8,["label"]),(0,l.Wm)(L.Z,{color:"primary",icon:"settings",label:e.$t("accountsUpdateQuotaDialog.updateBtn"),onClick:a[8]||(a[8]=e=>f())},null,8,["label"])])),_:1})])),_:1})])),_:1},8,["modelValue"])}}}),ga=ma,va=ga;v()(ma,"components",{QDialog:ee.Z,QCard:ae.Z,QBtn:te.Z,QCardSection:le.Z,QCardActions:oe.Z});var wa=t(7178);const ba=(0,l.aZ)({__name:"AccountsIndex",setup(e){const a=new M.Z,t=(0,N.o)(),s=(0,A.n)(),u=(0,o.iH)([]),n=(0,o.iH)(!1),c=(0,o.iH)(!1),i=(0,l.Fl)((()=>!0===n.value||!0===c.value)),r=(0,l.Fl)((()=>t.getKeyAccountsTable)),d=(0,l.Fl)((()=>t.getKeyUpdatePasswordDialog)),p=(0,l.Fl)((()=>t.getKeyUpdateApiKeyDialog)),m=(0,l.Fl)((()=>t.getKeyUpdateQuotaDialog)),g=(0,l.Fl)({get:()=>s.getSystemInfo,set:e=>s.setSystemInfo(e)});function v(){c.value=!0;const e=new wa.Z;e.getSystemInfo().then((e=>{g.value=e.data.body})).catch((e=>{console.error(e)})).finally((()=>{c.value=!1,w()}))}function w(){n.value=!0,a.getAccounts().then((e=>{u.value=e.data.body})).catch((e=>{console.error(e),(0,X.s9)(e.response.data,"accountsIndex.errorLoadingAccounts")})).finally((()=>{n.value=!1}))}return(0,l.wF)((()=>{v()})),(0,l.YP)(r,(()=>{w()})),(e,a)=>{const t=(0,l.up)("q-skeleton"),o=(0,l.up)("q-card-section"),s=(0,l.up)("q-card"),n=(0,l.up)("q-page");return(0,l.wg)(),(0,l.j4)(n,{padding:""},{default:(0,l.w5)((()=>[((0,l.wg)(),(0,l.j4)(Xe,{key:d.value})),((0,l.wg)(),(0,l.j4)(sa,{key:p.value})),(0,l.Wm)(ra),((0,l.wg)(),(0,l.j4)(va,{key:m.value})),(0,l.Wm)(s,{flat:""},{default:(0,l.w5)((()=>[(0,l.Wm)(o,null,{default:(0,l.w5)((()=>[!0===i.value?((0,l.wg)(),(0,l.j4)(t,{key:0,animation:"wave",style:{height:"100vh"}})):((0,l.wg)(),(0,l.j4)(Ve,{key:1,accountsList:u.value},null,8,["accountsList"]))])),_:1})])),_:1})])),_:1})}}});var ya=t(9885),fa=t(7133);const Qa=ba,Ua=Qa;v()(ba,"components",{QPage:ya.Z,QCard:ae.Z,QCardSection:le.Z,QSkeleton:fa.ZP})}}]);