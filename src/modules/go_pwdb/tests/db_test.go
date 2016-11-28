package tests

import (
	"config"
	"db"
	"fmt"
	"testing"
	"os"
	"sort"
)

func TestNoExtantDB( t *testing.T ){
	conf,err := config.CoreParser([]string{})
	if err != nil{
		fmt.Println("Failed to allocate config",err)
		t.FailNow()
	}
		
	conf.Path = "test.db"
	err = os.Remove(conf.Path)
	if err != nil && !os.IsNotExist(err){
		fmt.Println("Failed to clean out test db",err)
		t.FailNow()
	}
	
	err = db.Load(conf)
	if err != nil{
		fmt.Println("Failed to load config",err)
		t.FailNow()
	}
	db.Put("a_tag", "the_password")
	db.Put("another_tag", "password")
	pw,err := db.Get("a_tag")
	if err != nil{
		fmt.Println("Failed to fetch tag",err)
		t.FailNow()
	}
	if pw != "the_password" {
		fmt.Println("Unexpected password retrieved",err)
		t.FailNow()
	}
		
	err = db.Save(conf)
	if err != nil{
		fmt.Println("Failed to save config",err)
		t.FailNow()
	}
}

func TestExtantDB( t *testing.T ){
	conf,err := config.CoreParser([]string{})
	if err != nil{
		fmt.Println("Failed to allocate config",err)
		t.FailNow()
	}
		
	conf.Path = "test.db"
	prep_known_db_state( conf, t )
	
	err = db.Load(conf)
	if err != nil{
		fmt.Println("Failed to load config",err)
		t.FailNow()
	}
	pw,err := db.Get("a_tag")
	if err != nil{
		fmt.Println("Failed to fetch tag",err)
		t.FailNow()
	}
	if pw != "the_password" {
		fmt.Println("Unexpected password retrieved",err)
		t.FailNow()
	}
	
}

func TestGetAll( t *testing.T ){
	conf,err := config.CoreParser([]string{})
	if err != nil{
		fmt.Println("Failed to allocate config",err)
		t.FailNow()
	}
		
	conf.Path = "test.db"
	prep_known_db_state( conf, t )
	
	err = db.Load(conf)
	if err != nil{
		fmt.Println("Failed to load config",err)
		t.FailNow()
	}
	known_tags := db.GetAll()
	sort.Strings( known_tags )
	expect := []string{"a_tag", "another_tag"}
	sort.Strings( expect )
	fmt.Println("Expected tags ",expect)
	fmt.Println("Received tags ",known_tags)
	if known_tags[0] != expect[0] ||
		known_tags[1] != expect[1] {
		fmt.Println("Expected tags ",expect)
		fmt.Println("Received tags ",known_tags)
		t.FailNow()
	}
}		
	
func prep_known_db_state( conf config.Config, t *testing.T ) {
	err := os.Remove(conf.Path)
	if err != nil && !os.IsNotExist(err){
		fmt.Println("Failed to clean out old db",err)
		t.FailNow()
	}
	
	err = db.Load(conf)
	if err != nil{
		fmt.Println("Failed to load known state config",err)
		t.FailNow()
	}
	db.Put("a_tag", "the_password")
	db.Put("another_tag", "password")

	err = db.Save(conf)
	if err != nil{
		fmt.Println("Failed to save config",err)
		t.FailNow()
	}
}

