package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "ggnetwork"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}

}

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	// Root is the folder name of the root directory containing all the folders/files of the given file.
	Root string
	PathTransformFunc PathTransformFunc
	
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		Filename: key,
	}
}

type Store struct {
	StoreOpts
}

type PathKey struct {
	PathName string
	Filename string
}

func(p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName,"/")
	if len(paths)==0{
		return ""
	}
	return paths[0]
}


func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc==nil{
		opts.PathTransformFunc=DefaultPathTransformFunc
	}
	if len(opts.Root) == 0{
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	PathKey:=s.PathTransformFunc(key)
	fullPathWithRoot :=fmt.Sprintf("%s/%s",s.Root,PathKey.FullPath())
	_,err:=os.Stat(fullPathWithRoot)
	if errors.Is(err,os.ErrNotExist){
		return false
	}
	return true

}

func (s *Store) Delete (key string) error {
	pathKey:=s.PathTransformFunc(key)
	defer func(){
		log.Printf("deleted [%s] from disk",pathKey.Filename)
	}()
	 firstPathnameWithRoot:=fmt.Sprintf("%s/%s",s.Root,pathKey.FirstPathName())
	return os.RemoveAll(firstPathnameWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, nil

}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot :=fmt.Sprintf("%s/%s",s.Root,pathKey.FullPath())

	return os.Open(fullPathWithRoot)

}

func (s *Store) writestream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s",s.Root,pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPathWithRoot :=fmt.Sprintf("%s/%s",s.Root,pathKey.FullPath())
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Printf("written (%d) bytes to the disk: %s", n, fullPathWithRoot)
	return nil
}
