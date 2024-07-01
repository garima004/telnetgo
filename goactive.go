package main

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

func main() {
	// Replace these with appropriate values
	ldapServer := "ldap.example.com"
	ldapPort := 389
	bindUsername := "cn=admin,dc=example,dc=com"
	bindPassword := "adminpassword"
	baseDN := "dc=example,dc=com"

	// New user's details
	newUsername := "newuser"
	newPassword := "newuserpassword"
	newFirstName := "New"
	newLastName := "User"
	newUserDN := fmt.Sprintf("cn=%s,ou=Users,%s", newUsername, baseDN)

	// Connect to LDAP server
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// Bind with a read/write user
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		log.Fatal(err)
	}

	// Add new user
	addRequest := ldap.NewAddRequest(newUserDN, nil)
	addRequest.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})
	addRequest.Attribute("cn", []string{newUsername})
	addRequest.Attribute("sn", []string{newLastName})
	addRequest.Attribute("givenName", []string{newFirstName})
	addRequest.Attribute("displayName", []string{newFirstName + " " + newLastName})
	addRequest.Attribute("userPrincipalName", []string{newUsername + "@example.com"})
	addRequest.Attribute("sAMAccountName", []string{newUsername})
	addRequest.Attribute("userPassword", []string{newPassword})

	err = l.Add(addRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User %s added successfully.\n", newUsername)
}
