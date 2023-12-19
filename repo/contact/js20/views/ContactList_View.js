/** Copyright (c) 2022
 *	Andrey Mikhalevich, Katren ltd.
 */
function ContactList_View(id,options){	

	options = options || {};
	options.HEAD_TITLE = "Контакты";

	ContactList_View.superclass.constructor.call(this,id,options);
	
	var model = (options.models && options.models.ContactList_Model)? options.models.ContactList_Model : new ContactList_Model();
	var contr = new Contact_Controller();
	
	var constants = {"doc_per_page_count":null,"grid_refresh_interval":null};
	window.getApp().getConstantManager().get(constants);

	var popup_menu = new PopUpMenu();
	var pagClass = window.getApp().getPaginationClass();
	var self = this;
	this.addElement(new GridAjx(id+":grid",{
		"model":model,
		"controller":contr,
		"editInline":false,
		"editWinClass":Contact_Form,
		"commands":new GridCmdContainerAjx(id+":grid:cmd"),		
		"popUpMenu":popup_menu,
		"head":new GridHead(id+"-grid:head",{
			"elements":[
				new GridRow(id+":grid:head:row0",{
					"elements":[
						new GridCellHead(id+":grid:head:name",{
							"value":"Наименование",
							"columns":[
								new GridColumn({"field":model.getField("name")})
							],
							"sortable":true,
							"sort":"asc"							
						})
						,new GridCellHead(id+":grid:head:tm_first_name",{
							"value":"Telegram",
							"columns":[
								new GridColumn({
									"field":model.getField("tm_first_name"),
									"ctrlOptions":{
										"enabled":false,
										"labelCaption":""
									},
									"formatFunction":function(f, cell){
										var cell_n = cell.getNode();										
										var ft = f.tm_photo.getValue();
										if(ft){
											var i = document.createElement("img");
											i.setAttribute("src", "data:image/png;base64, "+ft);
											i.className = "userPhoto";
											cell_n.appendChild(i);
											if(self.photoDetail){
												delete self.photoDetail;
											}											
											self.photoDetail = new ToolTip({
												"node": cell_n,
												"wait":2,
												"onHover":(function(grid, id){
													return function(event){
														if(!grid.photoDetailData){
															grid.photoDetailData = [];
														}
														if(!grid.photoDetailData[id]){
															var pm = (new TmUser_Controller()).getPublicMethod("get_object");
															pm.setFieldValue("id", id);
															var ctrl = this;
															pm.run({
																"ok":function(resp){
																	var m = resp.getModel("TmUserDialog_Model");
																	if(m.getNextRow()){
																		self.photoDetailData[id] = m.getFieldValue("tm_photo");
																		self.showPhoto(ctrl, self.photoDetailData[id]);
																	}
																}																
															});
														}else{
															self.showPhoto(this, self.photoDetailData[id]);
														}													
													}
												})(self, f.ext_id.getValue())
											});
										}
										var nm = f.tm_first_name.getValue();	
										if(nm && nm.length){									
											var t = document.createElement("span");										
											DOMHelper.setText(t, nm);
											cell_n.appendChild(t);
										}
										return "";
									}
								})
							],
							"sortable":true
						})
						
						,new GridCellHead(id+":grid:head:email",{
							"value":"Email",
							"columns":[
								new GridColumn({
									"field":model.getField("email"),
									"ctrlClass":EditString,
									"ctrlOptions":{
										"maxLength":50
									}
								})
							],
							"sortable":true
						})
						,new GridCellHead(id+":grid:head:tel",{
							"value":"Телефон",
							"columns":[
								new GridColumnPhone({
									"field":model.getField("tel"),
									"ctrlClass":EditPhone
								})
							],
							"sortable":true
						})
						,new GridCellHead(id+":grid:head:tel_ext",{
							"value":"Добавочный",
							"columns":[
								new GridColumn({
									"field":model.getField("tel_ext")
								})
							]
						})
						
						,new GridCellHead(id+":grid:head:posts_ref",{
							"value":"Должность",
							"columns":[
								new GridColumnRef({
									"field":model.getField("posts_ref"),
									"ctrlClass":PostEdit,
									"ctrlOptions":{
										"labelCaption":""
									}
									
								})
							]
						})
						,new GridCellHead(id+":grid:head:comment_text",{
							"value":"Комментарий",
							"columns":[
								new GridColumn({
									"field":model.getField("comment_text")
								})
							]
						})
						
					]
				})
			]
		}),
		"pagination":new pagClass(id+"_page",
			{"countPerPage":constants.doc_per_page_count.getValue()}),		
		
		"autoRefresh":false,
		"refreshInterval":constants.grid_refresh_interval.getValue()*1000,
		"rowSelect":false,
		"focus":true
	}));	
	


}
extend(ContactList_View,ViewAjxList);

ContactList_View.prototype.showPhoto = function(ctrl, base64Data){
	ctrl.popup(
		'<div><img src="data:image/png;base64, '+base64Data+'"/></div>',
		{"title":"Данные по контакту"}
	);
}
