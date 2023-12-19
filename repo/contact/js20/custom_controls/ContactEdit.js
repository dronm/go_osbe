/** Copyright (c) 2019 
 *	Andrey Mikhalevich, Katren ltd.
 */
function ContactEdit(id,options){

	options = options || {};	
	if (options.labelCaption!=""){
		options.labelCaption = options.labelCaption || "Контакт:";
	}
	options.cmdInsert = (options.cmdInsert!=undefined)? options.cmdInsert:true;
	
	options.keyIds = options.keyIds || ["id"];
	
	//форма выбора из списка
	options.selectWinClass = ContactList_Form;
	options.selectDescrIds = options.selectDescrIds || ["descr"];
	
	//форма редактирования элемента
	options.editWinClass = Contact_Form;
	
	options.acMinLengthForQuery = (options.acMinLengthForQuery!=undefined)? options.acMinLengthForQuery:1;
	options.acController = new Contact_Controller();
	if(options.tm_activated == true){
		options.acPublicMethod = options.acController.getPublicMethod("complete_tm");
	}else{
		options.acPublicMethod = options.acController.getPublicMethod("complete");
	}
	options.acModel = new ContactList_Model();
	options.acPatternFieldId = options.acPatternFieldId || "descr";
	options.acKeyFields = options.acKeyFields || [options.acModel.getField("id")];
	options.acDescrFields = options.acDescrFields || [options.acModel.getField("descr")];
	options.acICase = options.acICase || "1";
	options.acMid = options.acMid || "1";
	options.acDescrFunction = function(f){
		var p_ref = f.posts_ref.getValue();
		//" +"+CommonHelper.maskFormat(f.tel.getValue(), window.getApp().getPhoneEditMask())
		return f["descr"].getValue() + ((p_ref && !p_ref.isNull())? " ("+p_ref.getDescr()+")": "");
	};	
	options.acOnCompleteTextOut = function(textHTML,modelRow){
		var pref = "";
		if(modelRow&&modelRow.tm_photo&&modelRow.tm_photo.getValue()){
			//Contact photo
			pref = "<img class='contactPhoto' src='data:image/png;base64, "+modelRow.tm_photo.getValue()+"'/img>";
		}
		return pref + textHTML;
	}
	
	ContactEdit.superclass.constructor.call(this,id,options);
}
extend(ContactEdit,EditRef);

/* Constants */


/* private members */

/* protected*/


/* public methods */

