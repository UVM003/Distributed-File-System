package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"testing"
)


func TestPathTRansformFunc(t *testing.T){
	key:= "momsbestpicture"
	passkey := CASPathTransformFunc(key)
	expectedOriginal:= "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathname:= "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if passkey.PathName!=expectedPathname{
		t.Errorf("have to get %s for %s", expectedPathname,passkey.PathName)
	}

	if passkey.Filename!=expectedOriginal{
		t.Errorf("have to get %s for %s", expectedOriginal,passkey.Filename)
	}
}

func TestStoreDeleteKey(t *testing.T){
	opts:= StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s:=NewStore(opts)
	key:="myspecial picture"
	data := []byte("Some jpg bytes")
	if err:=s.writestream(key,bytes.NewReader(data)); err !=nil{
		t.Error(err)
	}

	if err:= s.Delete(key); err !=nil{
		t.Error(err)
	}

	if s.Has(key){
		t.Errorf("The Delete is not implemented propely")
	}
 }


func TestStore(t *testing.T) {
	opts:=StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		
	}
	s:=NewStore(opts)
	key:="myspecial picture"
	data := []byte("Some jpg bytes")
	if err:=s.writestream(key,bytes.NewReader(data)); err !=nil{
		t.Error(err)
	}

	r,err:=s.readStream("myspecial picture")
	if err!=nil {
		 t.Error(err)
	}
	  
	b,_:=io.ReadAll(r)
  
	r.Close()
	fmt.Println(string(b))
	if string(b)!=string(data){
		t.Errorf("Wanted %s but got %s",data,b)
	}
	s.Delete(key)
	if s.Has(key){
		t.Errorf("The deletion was not successful")
	} else {
		log.Printf("The deletion was successful")
	}

}
