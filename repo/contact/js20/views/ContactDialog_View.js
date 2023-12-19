/**
 * @author Andrey Mikhalevich <katrenplus@mail.ru>, 2017
 
 * @extends ViewObjectAjx.js
 * @requires core/extend.js  
 * @requires controls/ViewObjectAjx.js 
 
 * @class
 * @classdesc
	
 * @param {string} id view identifier
 * @param {object} options
 * @param {object} options.models All data models
 * @param {object} options.variantStorage {name,model}
 */	
function ContactDialog_View(id,options){	

	options = options || {};
	
	options.controller = new Contact_Controller();
	options.model = (options.models&&options.models.ContactDialog_Model)? options.models.ContactDialog_Model : new ContactDialog_Model();
	
	var tm_inv_vis = true;
	if(options.model && options.model.getNextRow()){
		if(options.model.getFieldValue("tm_activated")){
			options.templateOptions = options.templateOptions || {};
			options.templateOptions.TM_ACTIVATED = true;
			options.templateOptions.TM_PHOTO = options.model.getFieldValue("tm_photo");
			tm_inv_vis = false;
		}
	}
	var self = this;
	options.addElement = function(){
		this.addElement(new EditString(id+":name",{
			"required": true,
			"focus": true,
			"maxLength":250,
			"labelCaption": "Фамилия имя отчество",
			"placeholder":"ФИО",
			"title":"Фамилия имя отчество физического лица"
		}));	

		this.addElement(new PostEdit(id+":posts_ref",{
			"placeholder": "Должность",
			"title":"Должность физического лица"
		}));	
			
		this.addElement(new EditEmail(id+":email",{
			"labelCaption":"Эл.почта:",
			"title":"Адрес электронной почты",
			"placeholder":"Адрес электронной почты"
		}));	

		this.addElement(new EditPhone(id+":tel",{
			"labelCaption":"Телефон:",
			"title":"Телефон физического лица"
		}));

		this.addElement(new EditString(id+":tel_ext",{
			"labelCaption":"Добавочный номер:",
			"placeholder":"Добавочный номер телефона",
			"title":"Добавочный номер телефона физического лиц"
		}));	
		this.addElement(new EditText(id+":comment_text",{
			"labelCaption":"Коммантарий:",
			"title":"Коммантарий менеджера о физическом лице"
		}));	

		/*
		this.addElement(new ButtonCtrl(id+":btnTmInvite",{
			"caption":"Пригласить в Telegram:",
			"title":"Отправить контакту приглашение в Telegram",
			"glyph":"glyphicon-send",
			"visible": tm_inv_vis,
			"onClick": function(){
				var b = this;
				b.setEnabled(false);
				window.getApp().TMInviteContact(
					new RefType({"keys" :{"id": self.getElement("id").getValue()},
						"descr": self.getElement("name").getValue(),
						"dataType": "contacts"
						}),
					function(){
						b.setEnabled(true);
					}
				);
			}
		}));	
		*/
	}
	
	ContactDialog_View.superclass.constructor.call(this,id,options);
	
	//****************************************************
	//read
	this.setDataBindings([
		new DataBinding({"control":this.getElement("name")})
		,new DataBinding({"control":this.getElement("posts_ref"), "fieldId":"posts_ref"})		
		,new DataBinding({"control":this.getElement("email")})		
		,new DataBinding({"control":this.getElement("tel")})
		,new DataBinding({"control":this.getElement("tel_ext")})
		,new DataBinding({"control":this.getElement("comment_text")})
	]);
	
	//write
	this.setWriteBindings([
		new CommandBinding({"control":this.getElement("name")})
		,new CommandBinding({"control":this.getElement("posts_ref"), "fieldId":"post_id"})		
		,new CommandBinding({"control":this.getElement("email")})
		,new CommandBinding({"control":this.getElement("tel")})
		,new CommandBinding({"control":this.getElement("tel_ext")})
		,new CommandBinding({"control":this.getElement("comment_text")})
	]);
}
extend(ContactDialog_View,ViewObjectAjx);
