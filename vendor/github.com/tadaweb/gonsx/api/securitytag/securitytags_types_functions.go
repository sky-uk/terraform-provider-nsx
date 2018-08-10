package securitytag

import "fmt"

func (s SecurityTags) String() string {
	var returnString = "Security tags contains a list of securitytags"
	return returnString
}

func (s SecurityTag) String() string {
	return fmt.Sprintf("Security tag name %s and id %s", s.Name, s.ObjectID)
}

// FilterByName - Filters the SecurityTags->SecurityTag with provided name
func (s SecurityTags) FilterByName(name string) *SecurityTag {
	var securityTagFound SecurityTag
	for _, securityTag := range s.SecurityTags {
		if securityTag.Name == name {
			securityTagFound = securityTag
			break
		}
	}
	return &securityTagFound
}

// CheckByName - Returns true or false depending if name is in securityTags
func (s SecurityTags) CheckByName(name string) bool {
	for _, securityTag := range s.SecurityTags {
		if securityTag.Name == name {
			return true
		}
	}
	return false
}

// FilterByIDAttached - Filters BasicInfoList->BasicInfo with provided name
func (b BasicInfoList) FilterByIDAttached(id string) *BasicInfo {
	var basicInfoFound BasicInfo
	for _, basicInfo := range b.BasicInfoList {
		if basicInfo.ObjectID == id {
			basicInfoFound = basicInfo
			break
		}
	}
	return &basicInfoFound
}

// AddSecurityTagToAttachmentList - Adds a SecurityTagAttachment to a SecurityTagAttachmentList
func (s *AttachmentList) AddSecurityTagToAttachmentList(st Attachment) {
	s.SecurityTagAttachments = append(s.SecurityTagAttachments, st)
}

// CheckByObjectID - Checks if a specific securityTagAttachment is within a SecurityTagAttachmentList
func (s *AttachmentList) CheckByObjectID(objectID string) bool {
	for _, securityTagAttachment := range s.SecurityTagAttachments {
		if securityTagAttachment.ObjectID == objectID {
			return true
		}
	}
	return false
}

// VerifyAttachments - Returns a list of tags that need to be detached from the VM
func (s *AttachmentList) VerifyAttachments(originalTags *SecurityTags) []string {
	var tagsToRemove []string
	for _, securityTag := range originalTags.SecurityTags {
		if s.CheckByObjectID(securityTag.ObjectID) == false {
			tagsToRemove = append(tagsToRemove, securityTag.ObjectID)
		}
	}
	return tagsToRemove
}
