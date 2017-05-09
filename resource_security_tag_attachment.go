package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitytag"
	"log"
	"errors"
)

func getAllSecurityTagsAttached(moid string, nsxclient *gonsx.NSXClient) (*securitytag.SecurityTags, error)  {
	getAllAttachedToVmAPI := securitytag.NewGetAllAttachedToVM(moid)
	err := nsxclient.Do(getAllAttachedToVmAPI)

	if err != nil {
		return nil, err
	}

	if getAllAttachedToVmAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAttachedToVmAPI.StatusCode(), getAllAttachedToVmAPI.ResponseObject())
	}
	securityTagsAttached := getAllAttachedToVmAPI.GetResponse()

	return  securityTagsAttached, err
}

func getAttachmentList(d *schema.ResourceData) (*securitytag.AttachmentList, error) {
	securityTags := new(securitytag.AttachmentList)
	log.Println("[DEBUG] Entered getAttachmentList")
	// Gather the attributes for the resource.
	if v, ok := d.GetOk("tagid"); ok {
		log.Println("[DEBUG] tagid is okay")
		idList := v.([]interface{})

		//securityTags := make([]securitytag.AttachmentList, len(idList))
		for _, value := range idList {
			log.Println("[DEBUG] Entered for loop")
			//*log.Println(value.(string))
			attachment := securitytag.Attachment{ObjectID: value.(string)}
			log.Println("[DEBUG] Attachment made from security tag")
			log.Println(attachment.ObjectID)
			securityTags.AddSecurityTagToAttachmentList(attachment)
			log.Println("[DEBUG] Added tag to attachment list")
			if !ok {
				log.Println("[DEBUG] empty element found in securityTags")
				return nil, fmt.Errorf("empty element found in securityTags")
			}
		}
	} else {
		return nil, fmt.Errorf("tagid argument is required")
	}
	log.Println(securityTags.SecurityTagAttachments)
	return securityTags, nil
}

func resourceSecurityTagAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityTagAttachmentCreate,
		Read:   resourceSecurityTagAttachmentRead,
		Delete: resourceSecurityTagAttachmentDelete,
		Update: resourceSecurityTagAttachmentUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:	schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tagid": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"moid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSecurityTagAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, moid string
	var tagIDs []string
	var securityTags *securitytag.AttachmentList

	// We get an AttachmentList from the resource data
	securityTags, tagError := getAttachmentList(d)
	if tagError != nil {
		return tagError
	}
	if v, ok := d.GetOk("tagid"); ok {
		tagList := v.([]interface{})
		tagIDs = make([]string, len(tagList))
		for i, value := range tagList {
			tagID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in securitygroups")
			}
			tagIDs[i] = tagID
		}
	} else {
		return fmt.Errorf("tagid argument is required")
	}

	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	createAPI := securitytag.NewUpdateAttachedTags(moid, securityTags)
	err := nsxclient.Do(createAPI)
	if err != nil {
		return err
	}

	if createAPI.StatusCode() != 200 {
		log.Printf(fmt.Sprintf("[DEBUG] Response %v", createAPI.ResponseObject()))
		return fmt.Errorf("Failed to attach security tag %s", tagIDs)
	}

	id := name + "/" + moid
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	if len(tagIDs) > 0 && moid != "" {
		d.SetId(id)
	} else {
		return errors.New("Can not establish the id of the created resource")
	}
	return resourceSecurityTagAttachmentRead(d, m)

}

func resourceSecurityTagAttachmentRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name,moid string
	var tagList []interface{}
	var tagIDs []string

	if v, ok := d.GetOk("tagid"); ok {
		tagList = v.([]interface{})
		tagIDs = make([]string, len(tagList))
		for i, value := range tagList {
			tagID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in securitytags")
			}
			tagIDs[i] = tagID
		}
	} else {
		return fmt.Errorf("tag argument is required")
	}

	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
		d.Set("moid", moid)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", moid))
	securityTagAttachedList, err := getAllSecurityTagsAttached(moid,nsxclient)


	log.Println("The tags attached to the vm are:")
	log.Println(securityTagAttachedList.SecurityTags)

	if err != nil {
		return err
	}

	id := name + "/" + moid
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	if len(tagIDs) > 0 && moid != "" {
		d.SetId(id)
	} else {
		return errors.New("Can not establish the id of the created resource")
	}

	return nil
}

func resourceSecurityTagAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var moid string
	var tagList []interface{}
	var tagIDs []string


	if v, ok := d.GetOk("tagid"); ok {
		tagList = v.([]interface{})
		tagIDs = make([]string, len(tagList))
		for i, value := range tagList {
			tagID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in securitytags")
			}
			tagIDs[i] = tagID
		}
	} else {
		return fmt.Errorf("tag argument is required")
	}

	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", moid))
	attachedTags, err := getAllSecurityTagsAttached(moid,nsxclient)

	if err != nil {
		return err
	}

	//If a tag has been manually removed from the VM, then we remove it from the list to detach
	for i, id := range tagIDs{
		found := false
		for _, securityTag := range attachedTags.SecurityTags {
			if securityTag.ObjectID == id {
				found = true
			}
		}
		if !found {
			tagIDs = append(tagIDs[:i],tagIDs[i+1:]...)
			log.Println("DEBUG: ID NO LONGER NEEDED TO DETACH:"+id)
		}
	}


	// If the resource has been removed manually, notify Terraform of this fact.
	if len(tagIDs) == 0 {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	for _, id := range tagIDs{
		detachAPI := securitytag.NewDetach(id, moid)
		err = nsxclient.Do(detachAPI)
		log.Println("DEBUG DETACHED TAG :"+id)
		if err != nil {
			return err
		}
	}

	// If we got here, the resource had existed, we deleted it and there was
	// no error.  Notify Terraform of this fact and return successful
	// completion.
	d.SetId("")

	for _, id := range tagIDs {
		log.Printf(fmt.Sprintf("[DEBUG] id: %s is detached from %s", id, moid))
	}

	return nil
}

func resourceSecurityTagAttachmentUpdate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, moid string
	var tagIDs []string
	var securityTags *securitytag.AttachmentList

	securityTags, tagError := getAttachmentList(d)

	if tagError != nil {
		return tagError
	}

	if v, ok := d.GetOk("tagid"); ok {
		tagList := v.([]interface{})
		tagIDs = make([]string, len(tagList))
		for i, value := range tagList {
			tagID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in securitygroups")
			}
			tagIDs[i] = tagID
		}
	} else {
		return fmt.Errorf("tagid argument is required")
	}

	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}


	attachedTags, err := getAllSecurityTagsAttached(moid,nsxclient)

	if err != nil {
		return err
	}

	// We check to see if any of the tag's currently attached to the VM need to be detached
	var tagsToDetach []string
	for _, tag := range attachedTags.SecurityTags {
		log.Println("DEBUG On Tag :"+tag.ObjectID)
		if !securityTags.CheckByObjectID(tag.ObjectID){
			log.Println("DEBUG Tag to detach :"+tag.ObjectID)
			tagsToDetach = append(tagsToDetach, tag.ObjectID)
		}
	}

	for _, id := range tagsToDetach{
		detachAPI := securitytag.NewDetach(id, moid)
		err = nsxclient.Do(detachAPI)
		log.Println("DEBUG DETACHED TAG :"+id)
		if err != nil {
			return err
		}
	}

	// We now attach the tag's included in the update payload

	for _, tag := range securityTags.SecurityTagAttachments{
		log.Println("Updating with tag :"+tag.ObjectID)
	}

	updateAPI := securitytag.NewUpdateAttachedTags(moid, securityTags)
	updateErr := nsxclient.Do(updateAPI)
	log.Println("UpdateAPI completed")
	if updateErr != nil {
		return err
	}

	if updateAPI.StatusCode() != 200 {
		log.Printf(fmt.Sprintf("[DEBUG] Response %v", updateAPI.ResponseObject()))
		return fmt.Errorf("Failed to attach security tags")
	}

	id := name + "/" + moid
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	if len(tagIDs) > 0 && moid != "" {
		d.SetId(id)
	} else {
		return errors.New("Can not establish the id of the updated resource")
	}
	return resourceSecurityTagAttachmentRead(d, m)
}


