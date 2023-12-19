/** Copyright (c) 2022
 *	Andrey Mikhalevich, Katren ltd.
 */
function EntityContactList_View(id,options){	

	options = options || {};
	options.HEAD_TITLE = "Контакты";

	EntityContactList_View.superclass.constructor.call(this,id,options);
	
	var model = (options.models && options.models.EntityContactList_Model)? options.models.EntityContactList_Model : new EntityContactList_Model();
	var contr = new EntityContact_Controller();
	
	var constants = {"doc_per_page_count":null,"grid_refresh_interval":null};
	window.getApp().getConstantManager().get(constants);

	var popup_menu = new PopUpMenu();
	var pagClass = window.getApp().getPaginationClass();
	var self = this;
	this.addElement(new GridAjx(id+":grid",{
		"model":model,
		"controller":contr,
		"editInline": true,
		"editWinClass": null,
		"commands":new GridCmdContainerAjx(id+":grid:cmd"),		
		"popUpMenu":popup_menu,
		"head":new GridHead(id+"-grid:head",{
			"elements":[
				new GridRow(id+":grid:head:row0",{
					"elements":[
						/*,new GridCellHead(id+":grid:head:entity_type",{
							"value":"Тип владельца",
							"columns":[
								new EnumGridColumn_entity_types({
									"field":model.getField("entity_type"),
									"ctrlEdit": false
								})
							]
						})
						,new GridCellHead(id+":grid:head:entities_ref",{
							"value":"Владелец",
							"columns":[
								new GridColumnRef({
									"field":model.getField("entities_ref"),
									"ctrlClass":EditCompound,
									"ctrlOptions":{
										"labelCaption":"",
										"possibleDataTypes":{
											"users":{
												"ctrlClass": UserEditRef,
												"ctrlOptions":{},
												"dataType": "users",
												"dataDescrLoc": "Пользователь"
											}
											,"clients":{
												"ctrlClass": ClientEdit,
												"ctrlOptions":{},
												"dataType": "clients",
												"dataDescrLoc": "Контрагет"
											}
											,"suppliers":{
												"ctrlClass": SupplierEdit,
												"ctrlOptions":{},
												"dataType": "suppliers",
												"dataDescrLoc": "Поставщик"
											}
											
										}
									},
									"ctrlBindFieldId": "entity_id"
								})
							]
						})
						*/
						new GridCellHead(id+":grid:head:contacts_ref",{
							"value":"Контакт",
							"columns":[
								new GridColumn({
									"field":model.getField("contacts_ref"),
									"cellOptions":{
										"title":"Двойной клик для редактивания контакта",
										"events":{
											/*Does not Work!!!
											"click": function(e){
												self.getElement("grid").onClick(e);
											},*/
											"dblclick":function(e){
												e.preventDefault();
												e.stopPropagation();
												self.openContact(e.target);
											}
										}
									},
									"ctrlClass":ContactEdit,
									"ctrlOptions":{
										"labelCaption":""
									},
									"formatFunction":function(f, cell){
										self.formatContact(f, cell);
										return "";
									},
									"ctrlBindFieldId": "contact_id"
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
		//"refreshInterval":constants.grid_refresh_interval.getValue()*1000,
		"rowSelect":false,
		"focus":true
	}));	
	


}
extend(EntityContactList_View,ViewAjxList);


EntityContactList_View.prototype.formatContact = function(f, cell){
	var attrs = f.contact_attrs.getValue();
	if(!attrs || !attrs["name"]){
		return;
	}
	
	var cell_n = cell.getNode();
	cell_n.setAttribute("contact_id", f.contact_id.getValue());
	var c_tag = document.createElement("DIV");
	
	var tel = attrs.tel;
	var tel_m = tel;
	if(tel_m && tel_m.length==10){
		tel_m = "+7"+tel;
	}
	else if(tel_m && tel_m.length==11){
		tel_m = "+7"+tel.substr(1);
	}
	var s_tag = document.createElement("DIV");
	s_tag.textContent = attrs["name"];
	s_tag.setAttribute("style","cursor:pointer;");
	if(attrs["email"]){
		s_tag.textContent+= ", "+attrs["email"];
	}
	if(attrs["post"]){
		s_tag.textContent+= ", "+attrs["post"];
	}		
	if(tel_m && tel_m.length){
		s_tag.textContent+= ", ";
		var t_tag = document.createElement("A");
		t_tag.setAttribute("href","tel:"+tel_m);
		t_tag.textContent = CommonHelper.maskFormat(tel, window.getApp().getPhoneEditMask());
		if(attrs["tel_ext"]){
			t_tag.textContent+= " ("+attrs["tel_ext"]+")";
		}		
		s_tag.appendChild(t_tag);
	}
	c_tag.appendChild(s_tag);
	cell_n.appendChild(c_tag);
	
}

EntityContactList_View.prototype.openContact = function(n){
	var td = DOMHelper.getParentByTagName(n, "TD");
	if(!td){
		return;
	}
	var id = td.getAttribute("contact_id");
	if(!id){
		return;
	}
	(new Contact_Form({
		"name": "Contact_Form_"+id,
		"keys": {"id": id},
		"params":{
			"cmd":"edit",
			"editViewOptions":{}
		}
	})).open();
	
}

