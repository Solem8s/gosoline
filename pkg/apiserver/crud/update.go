package crud

import (
	"context"
	"errors"
	"github.com/applike/gosoline/pkg/apiserver"
	"github.com/applike/gosoline/pkg/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

type updateHandler struct {
	transformer UpdateHandler
}

func NewUpdateHandler(transformer UpdateHandler) gin.HandlerFunc {
	uh := updateHandler{
		transformer: transformer,
	}

	return apiserver.CreateJsonHandler(uh)
}

func (uh updateHandler) GetInput() interface{} {
	return uh.transformer.GetUpdateInput()
}

func (uh updateHandler) Handle(ctx context.Context, request *apiserver.Request) (*apiserver.Response, error) {
	id, valid := apiserver.GetUintFromRequest(request, "id")

	if !valid {
		return nil, errors.New("no valid id provided")
	}

	repo := uh.transformer.GetRepository()
	model := uh.transformer.GetModel()
	err := repo.Read(ctx, id, model)

	if err != nil {
		return nil, err
	}

	err = uh.transformer.TransformUpdate(request.Body, model)

	if modelNotChanged(err) {
		return apiserver.NewStatusResponse(http.StatusNotModified), nil
	}

	if err != nil {
		return nil, err
	}

	err = repo.Update(ctx, model)

	exists := db.IsDuplicateEntryError(err)

	if exists {
		return apiserver.NewStatusResponse(http.StatusConflict), nil
	}

	if err != nil {
		return nil, err
	}

	reload := uh.transformer.GetModel()
	err = repo.Read(ctx, model.GetId(), reload)

	if err != nil {
		return nil, err
	}

	apiView := GetApiViewFromHeader(request.Header)
	out, err := uh.transformer.TransformOutput(reload, apiView)

	if err != nil {
		return nil, err
	}

	return apiserver.NewJsonResponse(out), nil
}

func modelNotChanged(err error) bool {
	return errors.Is(err, ErrModelNotChanged)
}
