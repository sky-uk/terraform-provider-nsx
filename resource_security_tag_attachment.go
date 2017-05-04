package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitytag"
	"log"
)

func getSingleSecurityTagAttached(tagid, moid string, nsxclient *gonsx.NSXClient) (*securitytag.BasicInfo, error) {

	// Gather all the resources that are associated with the specified
	// scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] securitytag.NewGetAllAttached()"))
	getAllAPI := securitytag.NewGetAllAttached(tagid)
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	securityTagAttached := getAllAPI.GetResponse().FilterByIDAttached(moid)

	return securityTagAttached, nil
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
			log.Println(value.(string))
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
	return securityTags, nil
}

func resourceSecurityTagAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityTagAttachmentCreate,
		Read:   resourceSecurityTagAttachmentRead,
		Delete: resourceSecurityTagAttachmentDelete,
		Update: resourceSecurityGroupUpdate,

		Schema: map[string]*schema.Schema{
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
	var tagid, moid string
	var securityTags *securitytag.AttachmentList

	securityTags, noTagError := getAttachmentList(d)
	if noTagError != nil {
		return noTagError
	}

	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitytag.NewAssign(%s, %s)", tagid, moid))
	//createAPI := securitytag.NewAssign(tagid, moid)
	createAPI := securitytag.NewUpdateAttachedTags(moid, securityTags)
	err := nsxclient.Do(createAPI)

	if err != nil {
		return err
	}

	if createAPI.StatusCode() != 200 {
		log.Printf(fmt.Sprintf("[DEBUG] Response %v", createAPI.ResponseObject()))
		return fmt.Errorf("Failed to attach security tag %s", tagid)
	}

	id := tagid + "/" + moid
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	if tagid != "" && moid != "" {
		d.SetId(id)
	} else {
		return errors.New("Can not establish the id of the created resource")
	}

	return resourceSecurityTagAttachmentRead(d, m)
}

func resourceSecurityTagAttachmentRead(d *schema.ResourceData, m interface{}) error {
	var moid string
	var tagid []string
	//var securityTags *securitytag.AttachmentList

	/*
	securityTags, noTagError := getAttachmentList(d)
	if noTagError != nil {
		return noTagError
	}
	*/
	if v, ok := d.GetOk("tagid"); ok {
		tagid = v.([]string)
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
	securityTagAttachedList := securitytag.NewGetAllAttachedToVM(moid)

	if securityTagAttachedList.Error() != nil {
		return securityTagAttachedList.Error()
	}

	//for _, value := range securityTags.SecurityTagAttachments  {
	for _, value := range tagid  {
		if securityTagAttachedList.GetResponse().CheckByName(value) == false {
			return fmt.Errorf("Not found %s", value)
		}
	}

	return nil
}

func resourceSecurityTagAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var tagid, moid string

	if v, ok := d.GetOk("tagid"); ok {
		tagid = v.(string)
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
	securityTagAttached, err := getSingleSecurityTagAttached(tagid, moid, nsxclient)
	id := securityTagAttached.ObjectID

	if err != nil {
		return err
	}
	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	detachAPI := securitytag.NewDetach(tagid, moid)
	err = nsxclient.Do(detachAPI)

	if err != nil {
		return err
	}

	// If we got here, the resource had existed, we deleted it and there was
	// no error.  Notify Terraform of this fact and return successful
	// completion.
	d.SetId("")
	log.Printf(fmt.Sprintf("[DEBUG] id %s detached.", id))

	return nil
}

func resourceSecurityTagAttachmentUpdate(d *schema.ResourceData, m interface{}) error {
	/*
		hasChanges := false

		nsxclient := m.(*gonsx.NSXClient)
		var vmid, moid string

		getAllSecurityTagsAttachedToVMAPI := securitytag.NewGetAllAttachedToVM(vmid)
		currentSecurityTags := getAllSecurityTagsAttachedToVMAPI.GetResponse().SecurityTags
	*/

	/*TODO
	Get all current tags assigned to the VM using NewGetAllAttachedToVM
	Check for if there are any tags currently attached which are not in the new payload of tags by using verifyAttachments method on new payload
	If the returned list is not empty, detach the tags within that list from the VM
	Use nsxclient to update with new payload
	*/

	return nil
}
