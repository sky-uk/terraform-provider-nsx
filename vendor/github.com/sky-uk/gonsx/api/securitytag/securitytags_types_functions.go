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
