"use strict";(globalThis["webpackChunkcontrol"]=globalThis["webpackChunkcontrol"]||[]).push([[742],{7089:(e,a,t)=>{t.r(a),t.d(a,{default:()=>ga});var l=t(9835),o=t(499),s=t(6970),u=t(1957),n=t(4304);const c={class:"q-pb-xs"},i={class:"row justify-between items-center"},d={class:"col-3"},r={class:"col-9 q-pl-md"},p=(0,l.aZ)({__name:"AccountsQuotaInput",props:{quota:{},quotaAvailable:{},title:{},inputLabel:{},step:{}},emits:["update:quota"],setup(e,{emit:a}){const t=e,u=a,p=(0,o.iH)(0),m=(0,o.iH)(0),v=(0,o.iH)(t.quota),g=(0,o.iH)(!1),w=(0,l.Fl)((()=>p.value<=t.quotaAvailable));return(0,l.YP)(v,(e=>{u("update:quota",e),!0!==g.value&&(p.value=e,m.value++)}),{immediate:!0}),(0,l.YP)(p,(e=>{g.value=!0,v.value=e,setTimeout((()=>{g.value=!1}),10)}),{immediate:!0}),(e,a)=>{var o;const u=(0,l.up)("q-slider");return(0,l.wg)(),(0,l.iD)(l.HY,null,[(0,l._)("div",c,(0,s.zw)(e.title),1),(0,l._)("div",i,[(0,l._)("div",d,[((0,l.wg)(),(0,l.j4)(n.Z,{label:e.inputLabel,value:p.value,"onUpdate:value":a[0]||(a[0]=e=>p.value=e),key:m.value,rules:[()=>w.value||e.$t("accountsQuotaInput.invalidValue")],suffix:"GB",maxLength:"5"},null,8,["label","value","rules"]))]),(0,l._)("div",r,[(0,l.Wm)(u,{modelValue:v.value,"onUpdate:modelValue":a[1]||(a[1]=e=>v.value=e),disable:0===e.quotaAvailable,min:0,max:e.quotaAvailable,step:null!==(o=t.step)&&void 0!==o?o:1,label:"","label-value":`${v.value} GB`,"switch-label-side":"","label-always":"",markers:""},null,8,["modelValue","disable","max","step","label-value"])])])],64)}}});var m=t(8423),v=t(9984),g=t.n(v);const w=p,y=w;g()(p,"components",{QSlider:m.Z});var b=t(9476);const f={class:"q-pb-xs"},Q={class:"row justify-between items-center"},U={class:"col-3"},D={class:"col-9 q-pl-md"},_=(0,l.aZ)({__name:"AccountsInodesQuota",props:{inodesQuota:{},inodesAvailable:{}},emits:["update:inodesQuota"],setup(e,{emit:a}){const t=e,u=a,c=(0,o.iH)(t.inodesQuota),i=(0,o.iH)(0),d=(0,o.iH)(0),r=(0,o.iH)(!1),p=(0,l.Fl)((()=>i.value<=t.inodesAvailable));return(0,l.YP)(c,(e=>{u("update:inodesQuota",e),!0!==r.value&&(i.value=e,d.value++)}),{immediate:!0}),(0,l.YP)(i,(()=>{r.value=!0,c.value=i.value,setTimeout((()=>{r.value=!1}),10)}),{immediate:!0}),(e,a)=>{const t=(0,l.up)("q-slider");return(0,l.wg)(),(0,l.iD)(l.HY,null,[(0,l._)("div",f,(0,s.zw)(e.$t("accountsUpdateQuotaDialog.labelInodesQuota")),1),(0,l._)("div",Q,[(0,l._)("div",U,[((0,l.wg)(),(0,l.j4)(n.Z,{label:e.$t("accountsUpdateQuotaDialog.labelInodesQuota"),value:i.value,"onUpdate:value":a[0]||(a[0]=e=>i.value=e),key:d.value,rules:[()=>p.value||e.$t("accountsUpdateQuotaDialog.invalidValue")]},null,8,["label","value","rules"]))]),(0,l._)("div",D,[(0,l.Wm)(t,{modelValue:c.value,"onUpdate:modelValue":a[1]||(a[1]=e=>c.value=e),disable:0===e.inodesAvailable,min:0,max:e.inodesAvailable,label:"","label-value":(0,o.SU)(b.Z)(c.value,!0,2),"switch-label-side":"","label-always":"",markers:""},null,8,["modelValue","disable","max","label-value"])])])],64)}}}),h=_,q=h;g()(_,"components",{QSlider:m.Z});var A=t(6397);const Z={class:"q-pb-lg"},x={class:"q-pb-xs"},W={class:"q-py-sm q-pr-sm"},k={class:"q-py-md q-pr-sm"},C={class:"q-pb-md q-pr-sm"},S=128,F=128,B=4096,P=1e8,H=(0,l.aZ)({__name:"AccountsQuotasForm",props:{cpuQuota:{},memoryQuota:{},storageQuota:{},inodesQuota:{}},emits:["update:cpuQuota","update:memoryQuota","update:storageQuota","update:inodesQuota"],setup(e,{emit:a}){const t=e,u=a,n=(0,A.n)(),c=(0,o.iH)(t.cpuQuota),i=(0,o.iH)(t.memoryQuota),d=(0,o.iH)(t.storageQuota),r=(0,o.iH)(t.inodesQuota),p=(0,l.Fl)((()=>n.getSystemInfo)),m=(0,l.Fl)((()=>{const e={freeBytes:0,totalInodes:0};return p.value.resourceUsage.storageInfo.forEach((a=>{"xfs"===a.fileSystem&&(e.freeBytes=B,e.totalInodes=P)})),e})),v=(0,l.Fl)((()=>({cpuCores:S,memoryBytes:F,diskBytes:m.value.freeBytes,inodes:m.value.totalInodes})));return(0,l.YP)((()=>t.cpuQuota),(e=>{c.value=e})),(0,l.YP)((()=>t.memoryQuota),(e=>{i.value=e})),(0,l.YP)((()=>t.storageQuota),(e=>{d.value=e})),(0,l.YP)((()=>t.inodesQuota),(e=>{r.value=e})),(0,l.YP)(c,(e=>{u("update:cpuQuota",e)})),(0,l.YP)(i,(e=>{u("update:memoryQuota",e)})),(0,l.YP)(d,(e=>{u("update:storageQuota",e)})),(0,l.YP)(r,(e=>{u("update:inodesQuota",e)})),(e,a)=>{const t=(0,l.up)("q-slider");return(0,l.wg)(),(0,l.iD)(l.HY,null,[(0,l._)("div",Z,[(0,l._)("div",x,(0,s.zw)(e.$t("accountsQuotasForm.cpuQuotaTitle")),1),(0,l.Wm)(t,{modelValue:c.value,"onUpdate:modelValue":a[0]||(a[0]=e=>c.value=e),disable:0===v.value.cpuCores,min:0,step:4,max:v.value.cpuCores,label:"","label-value":c.value+" Cores","switch-label-side":"","label-always":"",markers:"",class:"q-px-xs"},null,8,["modelValue","disable","max","label-value"])]),(0,l._)("div",W,[(0,l.Wm)(y,{quota:d.value,"onUpdate:quota":a[1]||(a[1]=e=>d.value=e),quotaAvailable:v.value.diskBytes,title:e.$t("accountsQuotasForm.storageQuotaTitle"),inputLabel:e.$t("accountsQuotasForm.storageQuotaInputLabel")},null,8,["quota","quotaAvailable","title","inputLabel"])]),(0,l._)("div",k,[(0,l.Wm)(y,{quota:i.value,"onUpdate:quota":a[2]||(a[2]=e=>i.value=e),quotaAvailable:v.value.memoryBytes,title:e.$t("accountsQuotasForm.memoryQuotaTitle"),step:8,inputLabel:e.$t("accountsQuotasForm.memoryQuotaInputLabel")},null,8,["quota","quotaAvailable","title","inputLabel"])]),(0,l._)("div",C,[(0,l.Wm)(q,{inodesQuota:r.value,"onUpdate:inodesQuota":a[3]||(a[3]=e=>r.value=e),inodesAvailable:v.value.inodes},null,8,["inodesQuota","inodesAvailable"])])],64)}}}),T=H,$=T;g()(H,"components",{QSlider:m.Z});var K=t(5786),z=t(2626),V=t(7713),I=t(3064),L=t(5971),j=t(8900),Y=t(3055),G=t(6647);const X={class:"flex justify-between items-center q-mb-sm"},M={class:"title-dialog"},N=(0,l.aZ)({__name:"AccountsCreateUserDialog",props:{showCreateUserDialog:{type:Boolean}},emits:["update:showCreateUserDialog"],setup(e,{emit:a}){const t=e,c=a,i=(0,G.QT)().t,d=(0,Y.o)(),r=(0,o.iH)(t.showCreateUserDialog),p=(0,o.iH)(""),m=(0,o.iH)(""),v=(0,o.iH)(0),g=(0,o.iH)(0),w=(0,o.iH)(0),y=(0,o.iH)(0),b=(0,o.iH)(!1),f=(0,l.Fl)((()=>p.value.length>=3)),Q=(0,l.Fl)({get:()=>d.getKeyAccountsTable,set:e=>{d.setKeyAccountsTable(e)}});function U(){r.value=!1}function D(){const e=new L.Z;e.createAccount({username:p.value,password:m.value,quota:{cpuCores:v.value,memoryBytes:g.value,diskBytes:w.value,inodes:y.value}}).then((()=>{(0,j.LX)(i("accountsCreateUserDialog.createdSuccessfully")),Q.value++,U()})).catch((e=>{(0,j.s9)(e.response.data,i("accountsCreateUserDialog.errorCreatingAccount"))}))}return(0,l.YP)((()=>t.showCreateUserDialog),(e=>{r.value=e})),(0,l.YP)(r,(e=>{c("update:showCreateUserDialog",e)})),(e,a)=>{const t=(0,l.up)("q-btn"),o=(0,l.up)("q-card-section"),c=(0,l.up)("q-card-actions"),i=(0,l.up)("q-card"),d=(0,l.up)("q-dialog");return(0,l.wg)(),(0,l.j4)(d,{modelValue:r.value,"onUpdate:modelValue":a[11]||(a[11]=e=>r.value=e),persistent:""},{default:(0,l.w5)((()=>[(0,l.Wm)(i,{class:"dialog-card-bg",style:{width:"700px","max-width":"80vw"}},{default:(0,l.w5)((()=>[(0,l._)("div",X,[(0,l._)("div",M,(0,s.zw)(e.$t("accountsCreateUserDialog.title")),1),(0,l.Wm)(t,{flat:"",round:"",dense:"",icon:"close",onClick:a[0]||(a[0]=e=>U())})]),(0,l.Wm)(o,{class:"q-px-xs"},{default:(0,l.w5)((()=>[(0,l.Wm)(n.Z,{label:e.$t("accountsCreateUserDialog.username"),icon:"person",value:p.value,"onUpdate:value":a[1]||(a[1]=e=>p.value=e),rules:[()=>""!==p.value||e.$t("accountsCreateUserDialog.usernameRequired"),()=>p.value.length>=3||e.$t("accountsCreateUserDialog.usernameMinLength")]},null,8,["label","value","rules"])])),_:1}),(0,l.Wm)(o,{class:"q-px-xs"},{default:(0,l.w5)((()=>[(0,l.Wm)(z.Z,{password:m.value,"onUpdate:password":a[2]||(a[2]=e=>m.value=e)},null,8,["password"]),(0,l.Wm)(V.Z,{password:m.value,"onUpdate:password":a[3]||(a[3]=e=>m.value=e),class:"float-right"},null,8,["password"])])),_:1}),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.wy)((0,l.Wm)(I.Z,{password:m.value,minPasswordLength:6,"onUpdate:isValidPassword":a[4]||(a[4]=e=>b.value=e)},null,8,["password"]),[[u.F8,m.value.length>0]])])),_:1}),(0,l.Wm)(o,{class:"q-px-md"},{default:(0,l.w5)((()=>[(0,l.Wm)($,{cpuQuota:v.value,"onUpdate:cpuQuota":a[5]||(a[5]=e=>v.value=e),memoryQuota:g.value,"onUpdate:memoryQuota":a[6]||(a[6]=e=>g.value=e),storageQuota:w.value,"onUpdate:storageQuota":a[7]||(a[7]=e=>w.value=e),inodesQuota:y.value,"onUpdate:inodesQuota":a[8]||(a[8]=e=>y.value=e)},null,8,["cpuQuota","memoryQuota","storageQuota","inodesQuota"])])),_:1}),(0,l.Wm)(c,{align:"between",class:"q-pt-md q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(K.Z,{label:e.$t("accountsCreateUserDialog.cancelBtn"),class:"q-mr-sm q-px-xl",color:"grey-8",onClick:a[9]||(a[9]=e=>U())},null,8,["label"]),(0,l.Wm)(K.Z,{color:"primary",icon:"person_add",label:e.$t("accountsCreateUserDialog.createBtn"),disable:!f.value||!b.value,onClick:a[10]||(a[10]=e=>D())},null,8,["label","disable"])])),_:1})])),_:1})])),_:1},8,["modelValue"])}}});var E=t(2074),R=t(4458),J=t(8879),O=t(3190),ee=t(1821);const ae=N,te=ae;g()(N,"components",{QDialog:E.Z,QCard:R.Z,QBtn:J.Z,QCardSection:O.Z,QCardActions:ee.Z});const le={class:"flex justify-end items-center"},oe={class:"td-action-bg-15"},se={class:"td-action-bg-20"},ue={class:"td-action-bg-25"},ne=(0,l.aZ)({__name:"AccountTableActions",props:{account:{}},setup(e){const a=e,t=(0,Y.o)(),o=(0,l.Fl)({get:()=>t.getKeyUpdatePasswordDialog,set:e=>{t.setKeyUpdatePasswordDialog(e)}}),u=(0,l.Fl)({get:()=>t.getKeyUpdateApiKeyDialog,set:e=>{t.setKeyUpdateApiKeyDialog(e)}}),n=(0,l.Fl)({get:()=>t.getKeyUpdateQuotaDialog,set:e=>{t.setKeyUpdateQuotaDialog(e)}});function c(){o.value++,t.setAccountSelected(a.account),t.setShowUpdatePasswordDialog(!0)}function i(){t.setAccountSelected(a.account),t.setShowDeleteUserDialog(!0)}function d(){u.value++,t.setAccountSelected(a.account),t.setShowUpdateApiKeyDialog(!0)}function r(){n.value++,t.setAccountSelected(a.account),t.setShowUpdateQuotaDialog(!0)}return(e,a)=>{const t=(0,l.up)("q-tooltip"),o=(0,l.up)("q-btn");return(0,l.wg)(),(0,l.iD)("div",le,[(0,l._)("div",oe,[(0,l.Wm)(o,{size:"md",icon:"lock",flat:"",onClick:a[0]||(a[0]=e=>c())},{default:(0,l.w5)((()=>[(0,l.Wm)(t,{class:"bg-primary text-white",style:{"font-size":"14px"},offset:[10,10]},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(e.$t("accountsTableActions.editPasswordBtn")),1)])),_:1})])),_:1})]),(0,l._)("div",se,[(0,l.Wm)(o,{size:"md",icon:"key",flat:"",onClick:a[1]||(a[1]=e=>d())},{default:(0,l.w5)((()=>[(0,l.Wm)(t,{class:"bg-primary text-white",style:{"font-size":"14px"},offset:[10,10]},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(e.$t("accountsTableActions.editApiKeyBtn")),1)])),_:1})])),_:1})]),(0,l._)("div",ue,[(0,l.Wm)(o,{size:"md",icon:"settings",flat:"",onClick:a[2]||(a[2]=e=>r())},{default:(0,l.w5)((()=>[(0,l.Wm)(t,{class:"bg-primary text-white",style:{"font-size":"14px"},offset:[10,10]},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(e.$t("accountsTableActions.editQuotaBtn")),1)])),_:1})])),_:1})]),(0,l._)("div",null,[(0,l.Wm)(o,{color:"negative",icon:"delete",size:"md",onClick:a[3]||(a[3]=e=>i())},{default:(0,l.w5)((()=>[(0,l.Wm)(t,{class:"bg-negative text-white",style:{"font-size":"14px"},offset:[10,10]},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(e.$t("accountsTableActions.deleteBtn")),1)])),_:1})])),_:1})])])}}});var ce=t(6858);const ie=ne,de=ie;g()(ne,"components",{QBtn:J.Z,QTooltip:ce.Z});var re=t(8428),pe=t(7250);const me=e=>((0,l.dD)("data-v-0776f328"),e=e(),(0,l.Cn)(),e),ve=me((()=>(0,l._)("th",{class:"text-center account-table-th"},"CPU",-1))),ge=me((()=>(0,l._)("th",{class:"text-center account-table-th"},"Memory",-1))),we=me((()=>(0,l._)("th",{class:"text-center account-table-th"},"Disk",-1))),ye=me((()=>(0,l._)("th",{class:"text-center account-table-th"},"Inodes",-1))),be=me((()=>(0,l._)("th",{class:"text-center account-table-th"},"UserId / GroupId",-1))),fe={style:{width:"150px"}},Qe={style:{width:"150px"}},Ue={style:{width:"150px"}},De={style:{width:"150px"}},_e={class:"text-center"},he={class:"row justify-center q-mt-md"},qe=(0,l.aZ)({__name:"AccountsTable",props:{accountsList:{}},setup(e){const a=e,t=(0,o.iH)(!1),u=(0,o.iH)(0),n=(0,o.iH)(""),c=(0,o.iH)([{name:"username",label:"Username",align:"left",field:"username",classes:"td-main-table",headerClasses:"bg-primary text-white",sortable:!0}]),i=(0,o.iH)({sortBy:"desc",descending:!1,page:1,rowsPerPage:10}),d=(0,l.Fl)((()=>Math.ceil(a.accountsList.length/i.value.rowsPerPage))),r=e=>parseFloat((0,pe.Z)(e,!1,2,"GiB")),p=e=>parseFloat((0,b.Z)(e,!1,1));function m(){u.value++,t.value=!0}return(e,a)=>{const v=(0,l.up)("q-icon"),g=(0,l.up)("q-input"),w=(0,l.up)("q-space"),y=(0,l.up)("q-th"),f=(0,l.up)("q-tr"),Q=(0,l.up)("q-td"),U=(0,l.up)("q-table"),D=(0,l.up)("q-pagination");return(0,l.wg)(),(0,l.iD)(l.HY,null,[e.accountsList.length>0?((0,l.wg)(),(0,l.j4)(te,{showCreateUserDialog:t.value,"onUpdate:showCreateUserDialog":a[0]||(a[0]=e=>t.value=e),key:u.value},null,8,["showCreateUserDialog"])):(0,l.kq)("",!0),(0,l.Wm)(U,{rows:e.accountsList,columns:c.value,filter:n.value,pagination:i.value,"onUpdate:pagination":a[3]||(a[3]=e=>i.value=e),"no-data-label":e.$t("accountsTable.noDataLabel"),"row-key":"key",color:"primary",flat:"",bordered:"","hide-pagination":""},{top:(0,l.w5)((()=>[(0,l.Wm)(g,{borderless:"",debounce:"300",color:"primary",modelValue:n.value,"onUpdate:modelValue":a[1]||(a[1]=e=>n.value=e),label:e.$t("accountsTable.searchInput")},{prepend:(0,l.w5)((()=>[(0,l.Wm)(v,{name:"search"})])),_:1},8,["modelValue","label"]),(0,l.Wm)(w),(0,l.Wm)(K.Z,{label:e.$t("accountsTable.createAccountBtn"),color:"primary",icon:"person_add",onClick:a[2]||(a[2]=e=>m())},null,8,["label"])])),header:(0,l.w5)((e=>[(0,l.Wm)(f,{props:e},{default:(0,l.w5)((()=>[((0,l.wg)(!0),(0,l.iD)(l.HY,null,(0,l.Ko)(e.cols,(a=>((0,l.wg)(),(0,l.j4)(y,{key:a.name,props:e,style:{"font-weight":"bold","font-size":"14px"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(a.label),1)])),_:2},1032,["props"])))),128)),ve,ge,we,ye,be,(0,l.Wm)(y)])),_:2},1032,["props"])])),body:(0,l.w5)((e=>[(0,l.Wm)(f,{props:e},{default:(0,l.w5)((()=>[((0,l.wg)(!0),(0,l.iD)(l.HY,null,(0,l.Ko)(e.cols,(a=>((0,l.wg)(),(0,l.j4)(Q,{key:a.name,props:e},{default:(0,l.w5)((()=>[(0,l.Uk)((0,s.zw)(a.value),1)])),_:2},1032,["props"])))),128)),(0,l._)("td",fe,[(0,l.Wm)(re.Z,{total:e.row.quota.cpuCores,progress:e.row.quotaUsage.cpuCores,suffix:1===e.row.quota.cpuCores?"Core":"Cores"},null,8,["total","progress","suffix"])]),(0,l._)("td",Qe,[(0,l.Wm)(re.Z,{total:r(e.row.quota.memoryBytes),progress:r(e.row.quotaUsage.memoryBytes),suffix:"GB"},null,8,["total","progress"])]),(0,l._)("td",Ue,[(0,l.Wm)(re.Z,{total:r(e.row.quota.diskBytes),progress:r(e.row.quotaUsage.diskBytes),suffix:"GB"},null,8,["total","progress"])]),(0,l._)("td",De,[(0,l.Wm)(re.Z,{total:p(e.row.quota.inodes),progress:p(e.row.quotaUsage.inodes),suffix:(0,o.SU)(b.s)(e.row.quota.inodes)},null,8,["total","progress","suffix"])]),(0,l._)("td",_e,(0,s.zw)(e.row.id)+" / "+(0,s.zw)(e.row.groupId),1),(0,l.Wm)(Q,{class:"text-right"},{default:(0,l.w5)((()=>[(0,l.Wm)(de,{account:e.row},null,8,["account"])])),_:2},1024)])),_:2},1032,["props"])])),_:1},8,["rows","columns","filter","pagination","no-data-label"]),(0,l._)("div",he,[(0,l.Wm)(D,{modelValue:i.value.page,"onUpdate:modelValue":a[4]||(a[4]=e=>i.value.page=e),color:"primary",max:d.value,size:"md"},null,8,["modelValue","max"])])],64)}}});var Ae=t(1639),Ze=t(422),xe=t(6611),We=t(2857),ke=t(136),Ce=t(1233),Se=t(1682),Fe=t(7220),Be=t(996);const Pe=(0,Ae.Z)(qe,[["__scopeId","data-v-0776f328"]]),He=Pe;g()(qe,"components",{QTable:Ze.Z,QInput:xe.Z,QIcon:We.Z,QSpace:ke.Z,QTr:Ce.Z,QTh:Se.Z,QTd:Fe.Z,QPagination:Be.Z});var Te=t(5273);const $e={class:"flex justify-between items-center"},Ke={class:"title-dialog"},ze=(0,l.aZ)({__name:"AccountsUpdatePasswordDialog",setup(e){const a=(0,G.QT)().t,t=(0,Y.o)(),n=(0,l.Fl)({get:()=>t.getShowUpdatePasswordDialog,set:e=>{t.setShowUpdatePasswordDialog(e)}}),c=(0,l.Fl)((()=>t.getAccountSelected.id)),i=(0,l.Fl)((()=>t.getAccountSelected.username)),d=(0,o.iH)(""),r=(0,o.iH)(!1),p=(0,l.Fl)({get:()=>t.keyAccountsTable,set:e=>{t.keyAccountsTable=e}});function m(){n.value=!1}function v(){(0,Te.Q)();const e=new L.Z;e.updateAccount({accountId:c.value,password:d.value}).then((()=>{(0,j.LX)(a("accountsUpdatePasswordDialog.updatedSuccessfully")),p.value++,m()})).catch((e=>{(0,j.s9)(e.response.data,a("accountsUpdatePasswordDialog.errorUpdatingAccount"))})).finally((()=>{(0,Te.Z)()}))}return(e,a)=>{const t=(0,l.up)("q-btn"),o=(0,l.up)("q-card-section"),c=(0,l.up)("q-card-actions"),p=(0,l.up)("q-card"),g=(0,l.up)("q-dialog");return(0,l.wg)(),(0,l.j4)(g,{modelValue:n.value,"onUpdate:modelValue":a[6]||(a[6]=e=>n.value=e),persistent:""},{default:(0,l.w5)((()=>[(0,l.Wm)(p,{style:{width:"700px","max-width":"80vw"},class:"dialog-card-bg"},{default:(0,l.w5)((()=>[(0,l._)("div",$e,[(0,l._)("div",Ke,(0,s.zw)(e.$t("accountsUpdatePasswordDialog.title",{username:i.value})),1),(0,l.Wm)(t,{flat:"",round:"",dense:"",icon:"close",onClick:a[0]||(a[0]=e=>m())})]),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(z.Z,{password:d.value,"onUpdate:password":a[1]||(a[1]=e=>d.value=e),minLength:6},null,8,["password"]),(0,l.Wm)(V.Z,{password:d.value,"onUpdate:password":a[2]||(a[2]=e=>d.value=e),class:"float-right"},null,8,["password"])])),_:1}),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.wy)((0,l.Wm)(I.Z,{minPasswordLength:6,password:d.value,"onUpdate:isValidPassword":a[3]||(a[3]=e=>r.value=e)},null,8,["password"]),[[u.F8,d.value.length>0]])])),_:1}),(0,l.Wm)(c,{align:"between",class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(K.Z,{label:e.$t("accountsUpdatePasswordDialog.cancelBtn"),color:"grey-8",onClick:a[4]||(a[4]=e=>m())},null,8,["label"]),(0,l.Wm)(K.Z,{color:"primary",icon:"key",label:e.$t("accountsUpdatePasswordDialog.updateBtn"),disable:!1===r.value,onClick:a[5]||(a[5]=e=>v())},null,8,["label","disable"])])),_:1})])),_:1})])),_:1},8,["modelValue"])}}}),Ve=ze,Ie=Ve;g()(ze,"components",{QDialog:E.Z,QCard:R.Z,QBtn:J.Z,QCardSection:O.Z,QCardActions:ee.Z});const Le={class:"flex justify-between items-center"},je={class:"title-dialog"},Ye={key:0},Ge={class:"box-api-key q-mb-sm"},Xe={class:"flex justify-between"},Me={class:"text-primary",style:{"font-size":"16px"}},Ne={style:{"font-size":"14px"}},Ee={key:1,style:{height:"50px"}},Re=(0,l.aZ)({__name:"AccountsUpdateApiKeyDialog",setup(e){const a=(0,G.QT)().t,t=(0,Y.o)(),u=(0,o.iH)(""),n=(0,l.Fl)({get:()=>t.getShowUpdateApiKeyDialog,set:e=>{t.setShowUpdateApiKeyDialog(e)}}),c=(0,l.Fl)((()=>t.getAccountSelected.id)),i=(0,l.Fl)((()=>t.getAccountSelected.username));function d(e){navigator.clipboard.writeText(e),(0,j.c0)({msg:`${a("accountsUpdateApiKeyDialog.copiedToClipboard")}`,position:"bottom",type:"primary",html:!0})}function r(){n.value=!1}function p(){(0,Te.Q)();const e=new L.Z;e.updateAccount({accountId:c.value,shouldUpdateApiKey:!0}).then((e=>{u.value=e.data.body,(0,j.LX)(a("accountsUpdateApiKeyDialog.updatedSuccessfullyWithApiKey"))})).catch((e=>{(0,j.s9)(e.response.data,a("accountsUpdateApiKeyDialog.errorUpdatingAccount"))})).finally((()=>{(0,Te.Z)()}))}return(e,a)=>{const t=(0,l.up)("q-btn"),o=(0,l.up)("q-icon"),c=(0,l.up)("q-card-section"),m=(0,l.up)("q-card-actions"),v=(0,l.up)("q-card"),g=(0,l.up)("q-dialog");return(0,l.wg)(),(0,l.j4)(g,{modelValue:n.value,"onUpdate:modelValue":a[4]||(a[4]=e=>n.value=e),persistent:""},{default:(0,l.w5)((()=>[(0,l.Wm)(v,{style:{width:"700px","max-width":"80vw"},class:"dialog-card-bg"},{default:(0,l.w5)((()=>[(0,l._)("div",Le,[(0,l._)("div",je,(0,s.zw)(e.$t("accountsUpdateApiKeyDialog.title",{username:i.value})),1),(0,l.Wm)(t,{flat:"",round:"",dense:"",icon:"close",onClick:a[0]||(a[0]=e=>r())})]),(0,l.Wm)(c,{class:"q-px-none"},{default:(0,l.w5)((()=>[u.value?((0,l.wg)(),(0,l.iD)("div",Ye,[(0,l._)("div",Ge,[(0,l._)("div",Xe,[(0,l._)("p",Me,(0,s.zw)(e.$t("accountsUpdateApiKeyDialog.generatedApiKey")),1),(0,l.Wm)(o,{name:"content_copy",size:"20px",class:"cursor-pointer icon-copy-api-key",onClick:a[1]||(a[1]=e=>d(u.value))})]),(0,l.Uk)(" "+(0,s.zw)(u.value),1)]),(0,l._)("small",Ne,[(0,l.Wm)(o,{color:"amber",size:"sm",name:"warning"}),(0,l.Uk)(" "+(0,s.zw)(e.$t("accountsUpdateApiKeyDialog.saveNewApiKey")),1)])])):((0,l.wg)(),(0,l.iD)("div",Ee))])),_:1}),(0,l.Wm)(m,{align:"between",class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(K.Z,{label:e.$t("accountsUpdateApiKeyDialog.closeDialogBtn"),color:"grey-8",onClick:a[2]||(a[2]=e=>r())},null,8,["label"]),(0,l.Wm)(K.Z,{icon:"key",label:e.$t("accountsUpdateApiKeyDialog.generateNewApiKeyBtn"),color:"primary",onClick:a[3]||(a[3]=e=>p())},null,8,["label"])])),_:1})])),_:1})])),_:1},8,["modelValue"])}}}),Je=Re,Oe=Je;g()(Re,"components",{QDialog:E.Z,QCard:R.Z,QBtn:J.Z,QCardSection:O.Z,QIcon:We.Z,QCardActions:ee.Z});var ea=t(9906),aa=t(9302);const ta=(0,l.aZ)({__name:"AccountsDeleteUserDialog",setup(e){const a=(0,aa.Z)(),t=(0,G.QT)().t,o=(0,Y.o)(),s=(0,l.Fl)((()=>a.dark.isActive?"/_/icons/bomb_dark.svg":"/_/icons/bomb_light.svg")),u=(0,l.Fl)({get:()=>o.getShowDeleteUserDialog,set:e=>{o.setShowDeleteUserDialog(e)}}),n=(0,l.Fl)((()=>o.getAccountSelected.username)),c=(0,l.Fl)((()=>o.getAccountSelected.id)),i=(0,l.Fl)({get:()=>o.keyAccountsTable,set:e=>{o.keyAccountsTable=e}});function d(){u.value=!1}function r(){(0,Te.Q)();const e=new L.Z;e.deleteAccount(c.value).then((()=>{(0,j.LX)(t("accountsDeleteUserDialog.deletedSuccessfully",{username:n.value})),i.value++,d()})).catch((e=>{(0,j.s9)(e.response.data,t("accountsDeleteUserDialog.errorDeletingAccount",{username:n.value}))})).finally((()=>{(0,Te.Z)()}))}return(e,a)=>((0,l.wg)(),(0,l.j4)(ea.Z,{showDeleteDialog:u.value,"onUpdate:showDeleteDialog":a[2]||(a[2]=e=>u.value=e),titleDialog:e.$t("accountsDeleteUserDialog.title",{username:n.value}),imagePath:s.value,messageToDelete:e.$t("accountsDeleteUserDialog.messageDeleteAccount",{username:n.value}),warningToDelete:e.$t("accountsDeleteUserDialog.warningDeleteAccount")},{"card-actions":(0,l.w5)((()=>[(0,l.Wm)(K.Z,{label:e.$t("accountsDeleteUserDialog.cancelBtn"),color:"grey-8",onClick:a[0]||(a[0]=e=>d())},null,8,["label"]),(0,l.Wm)(K.Z,{color:"negative",label:e.$t("accountsDeleteUserDialog.deleteBtn"),onClick:a[1]||(a[1]=e=>r())},null,8,["label"])])),_:1},8,["showDeleteDialog","titleDialog","imagePath","messageToDelete","warningToDelete"]))}}),la=ta,oa=la,sa={class:"flex justify-between items-center"},ua={class:"title-dialog"},na=(0,l.aZ)({__name:"AccountsUpdateQuotaDialog",setup(e){const a=(0,G.QT)().t,t=(0,Y.o)(),u=(0,o.iH)(0),c=(0,o.iH)(0),i=(0,o.iH)(0),d=(0,o.iH)(0),r=(0,l.Fl)((()=>t.getAccountSelected.id)),p=(0,l.Fl)((()=>t.getAccountSelected)),m=(0,l.Fl)({get:()=>t.getKeyAccountsTable,set:e=>{t.setKeyAccountsTable(e)}}),v=(0,l.Fl)({get:()=>t.getShowUpdateQuotaDialog,set:e=>{t.setShowUpdateQuotaDialog(e)}}),g=(0,l.Fl)((()=>t.getAccountSelected.username));function w(){v.value=!1}function y(){const e=new L.Z;e.updateAccount({accountId:r.value,quota:{cpuCores:u.value,memoryBytes:c.value,diskBytes:i.value,inodes:d.value}}).then((()=>{(0,j.LX)(a("accountsUpdateQuotaDialog.updatedSuccessfully")),m.value++,w()})).catch((e=>{(0,j.s9)(e.response.data,a("accountsUpdateQuotaDialog.errorUpdatingQuota"))}))}return(0,l.YP)(v,(e=>{!1!==e&&(u.value=p.value.quota.cpuCores,c.value=p.value.quota.memoryBytes,i.value=p.value.quota.diskBytes,d.value=p.value.quota.inodes)}),{immediate:!0}),(e,a)=>{const t=(0,l.up)("q-btn"),o=(0,l.up)("q-card-section"),r=(0,l.up)("q-card-actions"),p=(0,l.up)("q-card"),m=(0,l.up)("q-dialog");return(0,l.wg)(),(0,l.j4)(m,{modelValue:v.value,"onUpdate:modelValue":a[8]||(a[8]=e=>v.value=e),persistent:""},{default:(0,l.w5)((()=>[(0,l.Wm)(p,{class:"dialog-card-bg",style:{width:"700px","max-width":"80vw","overflow-x":"hidden"}},{default:(0,l.w5)((()=>[(0,l._)("div",sa,[(0,l._)("div",ua,(0,s.zw)(e.$t("accountsUpdateQuotaDialog.title",{username:g.value})),1),(0,l.Wm)(t,{flat:"",round:"",dense:"",icon:"close",onClick:a[0]||(a[0]=e=>w())})]),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(n.Z,{label:e.$t("accountsUpdateQuotaDialog.username"),disable:"",icon:"person",value:g.value,"onUpdate:value":a[1]||(a[1]=e=>g.value=e)},null,8,["label","value"])])),_:1}),(0,l.Wm)(o,{class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)($,{cpuQuota:u.value,"onUpdate:cpuQuota":a[2]||(a[2]=e=>u.value=e),memoryQuota:c.value,"onUpdate:memoryQuota":a[3]||(a[3]=e=>c.value=e),storageQuota:i.value,"onUpdate:storageQuota":a[4]||(a[4]=e=>i.value=e),inodesQuota:d.value,"onUpdate:inodesQuota":a[5]||(a[5]=e=>d.value=e)},null,8,["cpuQuota","memoryQuota","storageQuota","inodesQuota"])])),_:1}),(0,l.Wm)(r,{align:"between",class:"q-px-none"},{default:(0,l.w5)((()=>[(0,l.Wm)(K.Z,{label:e.$t("accountsUpdateQuotaDialog.cancelBtn"),class:"q-mr-sm q-px-xl",color:"grey-8",onClick:a[6]||(a[6]=e=>w())},null,8,["label"]),(0,l.Wm)(K.Z,{color:"primary",icon:"settings",label:e.$t("accountsUpdateQuotaDialog.updateBtn"),onClick:a[7]||(a[7]=e=>y())},null,8,["label"])])),_:1})])),_:1})])),_:1},8,["modelValue"])}}}),ca=na,ia=ca;g()(na,"components",{QDialog:E.Z,QCard:R.Z,QBtn:J.Z,QCardSection:O.Z,QCardActions:ee.Z});var da=t(7178);const ra=(0,l.aZ)({__name:"AccountsIndex",setup(e){const a=new L.Z,t=(0,Y.o)(),s=(0,A.n)(),u=(0,o.iH)([]),n=(0,o.iH)(!1),c=(0,o.iH)(!1),i=(0,l.Fl)((()=>!0===n.value||!0===c.value)),d=(0,l.Fl)((()=>t.getKeyAccountsTable)),r=(0,l.Fl)((()=>t.getKeyUpdatePasswordDialog)),p=(0,l.Fl)((()=>t.getKeyUpdateApiKeyDialog)),m=(0,l.Fl)((()=>t.getKeyUpdateQuotaDialog)),v=(0,l.Fl)({get:()=>s.getSystemInfo,set:e=>s.setSystemInfo(e)});function g(){c.value=!0;const e=new da.Z;e.getSystemInfo().then((e=>{v.value=e.data.body})).catch((e=>{console.error(e)})).finally((()=>{c.value=!1,w()}))}function w(){n.value=!0,a.getAccounts().then((e=>{u.value=e.data.body})).catch((e=>{console.error(e),(0,j.s9)(e.response.data,"accountsIndex.errorLoadingAccounts")})).finally((()=>{n.value=!1}))}return(0,l.wF)((()=>{g()})),(0,l.YP)(d,(()=>{w()})),(e,a)=>{const t=(0,l.up)("q-skeleton"),o=(0,l.up)("q-card-section"),s=(0,l.up)("q-card"),n=(0,l.up)("q-page");return(0,l.wg)(),(0,l.j4)(n,{padding:""},{default:(0,l.w5)((()=>[((0,l.wg)(),(0,l.j4)(Ie,{key:r.value})),((0,l.wg)(),(0,l.j4)(Oe,{key:p.value})),(0,l.Wm)(oa),((0,l.wg)(),(0,l.j4)(ia,{key:m.value})),(0,l.Wm)(s,{flat:""},{default:(0,l.w5)((()=>[(0,l.Wm)(o,null,{default:(0,l.w5)((()=>[!0===i.value?((0,l.wg)(),(0,l.j4)(t,{key:0,animation:"wave",style:{height:"100vh"}})):((0,l.wg)(),(0,l.j4)(He,{key:1,accountsList:u.value},null,8,["accountsList"]))])),_:1})])),_:1})])),_:1})}}});var pa=t(9885),ma=t(7133);const va=ra,ga=va;g()(ra,"components",{QPage:pa.Z,QCard:R.Z,QCardSection:O.Z,QSkeleton:ma.ZP})}}]);