/**	
 *
 * THIS FILE IS GENERATED FROM TEMPLATE build/templates/models/Model_js.xsl
 * ALL DIRECT MODIFICATIONS WILL BE LOST WITH THE NEXT BUILD PROCESS!!!
 *
 * @author Andrey Mikhalevich <katrenplus@mail.ru>, 2017
 * @class
 * @classdesc Model class. Created from template build/templates/models/Model_js.xsl. !!!DO NOT MODEFY!!!
 
 * @extends ModelXML
 
 * @requires core/extend.js
 * @requires core/ModelXML.js
 
 * @param {string} id 
 * @param {Object} options
 */

function Contact_Model(options){
	var id = 'Contact_Model';
	options = options || {};
	
	options.fields = {};
	
			
				
			
				
	
	var filed_options = {};
	filed_options.primaryKey = true;	
	
	filed_options.autoInc = true;	
	
	options.fields.id = new FieldInt("id",filed_options);
	
				
	
	var filed_options = {};
	filed_options.primaryKey = false;	
	filed_options.alias = 'Наименование';
	filed_options.autoInc = false;	
	
	options.fields.name = new FieldString("name",filed_options);
	options.fields.name.getValidator().setRequired(true);
	options.fields.name.getValidator().setMaxLength('250');
	
				
	
	var filed_options = {};
	filed_options.primaryKey = false;	
	filed_options.alias = 'Должность';
	filed_options.autoInc = false;	
	
	options.fields.post_id = new FieldInt("post_id",filed_options);
	
				
	
	var filed_options = {};
	filed_options.primaryKey = false;	
	filed_options.alias = 'Email';
	filed_options.autoInc = false;	
	
	options.fields.email = new FieldString("email",filed_options);
	options.fields.email.getValidator().setMaxLength('100');
	
				
	
	var filed_options = {};
	filed_options.primaryKey = false;	
	filed_options.alias = 'Телефон';
	filed_options.autoInc = false;	
	
	options.fields.tel = new FieldString("tel",filed_options);
	options.fields.tel.getValidator().setMaxLength('11');
	
				
	
	var filed_options = {};
	filed_options.primaryKey = false;	
	filed_options.alias = 'Добавочный номер';
	filed_options.autoInc = false;	
	
	options.fields.tel_ext = new FieldString("tel_ext",filed_options);
	options.fields.tel_ext.getValidator().setMaxLength('20');
	
				
	
	var filed_options = {};
	filed_options.primaryKey = false;	
	filed_options.alias = 'Описание для поиска';
	filed_options.autoInc = false;	
	
	options.fields.descr = new FieldText("descr",filed_options);
	
				
	
	var filed_options = {};
	filed_options.primaryKey = false;	
	filed_options.alias = 'Комментарий';
	filed_options.autoInc = false;	
	
	options.fields.comment_text = new FieldText("comment_text",filed_options);
	
				
	
	var filed_options = {};
	filed_options.primaryKey = false;	
	filed_options.defValue = true;
	filed_options.alias = 'Адрес электр.почты подтвержден';
	filed_options.autoInc = false;	
	
	options.fields.email_confirmed = new FieldBool("email_confirmed",filed_options);
	
			
						
		Contact_Model.superclass.constructor.call(this,id,options);
}
extend(Contact_Model,ModelXML);

