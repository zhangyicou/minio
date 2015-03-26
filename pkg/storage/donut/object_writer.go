package donut

import (
	"encoding/json"
	"github.com/minio-io/iodine"
	"io/ioutil"
	"os"
	"path"
)

func newDonutObjectWriter(objectDir string) (Writer, error) {
	dataFile, err := os.OpenFile(path.Join(objectDir, "data"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, iodine.Error(err, map[string]string{"objectDir": objectDir})
	}
	return donutObjectWriter{
		root:          objectDir,
		file:          dataFile,
		metadata:      make(map[string]string),
		donutMetadata: make(map[string]string),
	}, nil
}

type donutObjectWriter struct {
	root          string
	file          *os.File
	metadata      map[string]string
	donutMetadata map[string]string
	err           error
}

func (d donutObjectWriter) Write(data []byte) (int, error) {
	written, err := d.file.Write(data)
	return written, iodine.Error(err, nil)
}

func (d donutObjectWriter) Close() error {
	if d.err != nil {
		return iodine.Error(d.err, nil)
	}
	metadata, _ := json.Marshal(d.metadata)
	ioutil.WriteFile(path.Join(d.root, "metadata.json"), metadata, 0600)
	donutMetadata, _ := json.Marshal(d.donutMetadata)
	ioutil.WriteFile(path.Join(d.root, "donutMetadata.json"), donutMetadata, 0600)

	return iodine.Error(d.file.Close(), nil)
}

func (d donutObjectWriter) CloseWithError(err error) error {
	if d.err != nil {
		d.err = err
	}
	return iodine.Error(d.Close(), nil)
}

func (d donutObjectWriter) SetMetadata(metadata map[string]string) error {
	for k := range d.metadata {
		delete(d.metadata, k)
	}
	for k, v := range metadata {
		d.metadata[k] = v
	}
	return nil
}

func (d donutObjectWriter) GetMetadata() (map[string]string, error) {
	metadata := make(map[string]string)
	for k, v := range d.metadata {
		metadata[k] = v
	}
	return metadata, nil
}

func (d donutObjectWriter) SetDonutMetadata(metadata map[string]string) error {
	for k := range d.donutMetadata {
		delete(d.donutMetadata, k)
	}
	for k, v := range metadata {
		d.donutMetadata[k] = v
	}
	return nil
}

func (d donutObjectWriter) GetDonutMetadata() (map[string]string, error) {
	donutMetadata := make(map[string]string)
	for k, v := range d.donutMetadata {
		donutMetadata[k] = v
	}
	return donutMetadata, nil
}