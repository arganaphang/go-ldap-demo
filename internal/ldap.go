package internal

import (
	"errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type UserLDAPData struct {
	ID       string
	Email    string
	Name     string
	FullName string
}

func AuthUsingLDAP(username, password string) (bool, *UserLDAPData, error) {
	l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%s", LDAP_HOST, LDAP_PORT))
	if err != nil {
		return false, nil, err
	}
	defer l.Close()
	err = l.Bind(LDAP_BIND_DN, LDAP_BIND_PASSWORD)
	if err != nil {
		return false, nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		LDAP_SEARCH_DN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn", "cn", "sn", "mail"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) == 0 {
		return false, nil, errors.New("user not found")
	}
	entry := sr.Entries[0]

	err = l.Bind(entry.DN, password)
	if err != nil {
		return false, nil, err
	}

	data := new(UserLDAPData)
	data.ID = username

	for _, attr := range entry.Attributes {
		switch attr.Name {
		case "sn":
			data.Name = attr.Values[0]
		case "mail":
			data.Email = attr.Values[0]
		case "cn":
			data.FullName = attr.Values[0]
		}
	}

	return true, data, nil
}

func AuthUsingLDAPPasswordless(username string) (bool, *UserLDAPData, error) {
	l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%s", LDAP_HOST, LDAP_PORT))
	if err != nil {
		return false, nil, err
	}
	defer l.Close()
	err = l.Bind(LDAP_BIND_DN, LDAP_BIND_PASSWORD)
	if err != nil {
		return false, nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		LDAP_SEARCH_DN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn", "cn", "sn", "mail"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) == 0 {
		return false, nil, errors.New("user not found")
	}
	entry := sr.Entries[0]

	err = l.UnauthenticatedBind(entry.DN)
	if err != nil {
		return false, nil, err
	}

	data := new(UserLDAPData)
	data.ID = username

	for _, attr := range entry.Attributes {
		switch attr.Name {
		case "sn":
			data.Name = attr.Values[0]
		case "mail":
			data.Email = attr.Values[0]
		case "cn":
			data.FullName = attr.Values[0]
		}
	}

	return true, data, nil
}
