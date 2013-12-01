package v2

import (
	"net/http"
	"strconv"

	etcdErr "github.com/coreos/etcd/error"
	"github.com/gorilla/mux"
)

func DeleteHandler(w http.ResponseWriter, req *http.Request, s Server) error {
	vars := mux.Vars(req)
	key := "/" + vars["key"]

	req.ParseForm()
	_, valueOk := req.Form["prevValue"]
	_, indexOk := req.Form["prevIndex"]
	recursive := (req.Form.Get("recursive") == "true")

	if !valueOk && !indexOk {
		c := s.Store().CommandFactory().CreateDeleteCommand(key, recursive)
		return s.Dispatch(c, w, req)
	}

	var err error
	prevIndex := uint64(0)
	prevValue := req.Form.Get("prevValue")

	if indexOk {
		prevIndexStr := req.Form.Get("prevIndex")
		prevIndex, err = strconv.ParseUint(prevIndexStr, 10, 64)

		// bad previous index
		if err != nil {
			return etcdErr.NewError(etcdErr.EcodeIndexNaN, "CompareAndDelete", s.Store().Index())
		}
	}

	if valueOk {
		if prevValue == "" {
			return etcdErr.NewError(etcdErr.EcodePrevValueRequired, "CompareAndDelete", s.Store().Index())
		}
	}

	c := s.Store().CommandFactory().CreateCompareAndDeleteCommand(key, recursive, prevValue, prevIndex)
	return s.Dispatch(c, w, req)
}
