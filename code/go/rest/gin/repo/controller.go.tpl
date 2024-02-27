package controllers

  import (
    "github.com/gin-gonic/gin"
    "github.com/pkg/errors"
    "net/http"
    "stages_user_api/core"
    "stages_user_api/core/query"
    "stages_user_api/ports/api/http/controllers/utils"
  )

  func New{{ .Name | singular | pascal }}Controller(app *core.Application) *{{ .Name | singular | pascal }}Controller {
    return &{{ .Name | singular | pascal }}Controller{app: app}
  }

  type {{ .Name | singular | pascal }}Controller struct {
    app *core.Application
  }

  {{ template "go-rest-blocks-swag-get" . }}
  func (controller *{{ .Name | singular | pascal }}Controller) Get(c *gin.Context)  {
    id, err := strconv.Atoi(c.Param("id"))

    if err != nil {
      c.JSON(http.StatusBadRequest, err)
      return
    }

    items, err := controller.app.{{ .Name | singular | pascal}}Repository.Get(c.Request.Context(), id)

    if err != nil {
      _ = c.Error(err)
      return
    }
    
    c.JSON(http.StatusOK, items)
  }

  {{ template "go-rest-blocks-swag-post" . }}
  func (controller *{{ .Name | singular | pascal }}Controller) Create(c *gin.Context)  {
    m := &models.TaskSimulation{}

    err := c.ShouldBind(m)
    if err != nil {
      _ = c.Error(core.NewRequestValidationError(err))
      return
    }

    id, err := controller.app.{{ .Name | singular | pascal}}Repository.Create(c.Request.Context(), m)
    if err != nil {
      _ = c.Error(err)
      return
    }

    item, err := controller.app.{{ .Name | singular | pascal}}Repository.Get(c.Request.Context(), id)
    if err != nil {
      _ = c.Error(err)
      return
    }

    c.JSON(http.StatusCreated, item)
  }

  {{ template "go-rest-blocks-swag-delete" . }}
  func (controller *{{ .Name | singular | pascal }}Controller) Update(c *gin.Context)  {
    
    id, err := strconv.Atoi(c.Param("id"))

    if err != nil {
      c.JSON(http.StatusBadRequest, err)
      return
    }

    err = controller.app.{{ .Name | singular | pascal}}Repository.Delete(c.Request.Context(), id)
    if err != nil {
      if err == core.ErrRowNotFound {
        c.JSON(http.StatusNotFound, err)
        return
      }
      _ = c.Error(err)
      return
    }

    c.Status(http.StatusNoContent)
  }

  {{ template "go-rest-blocks-swag-put" . }}
  func (controller *{{ .Name | singular | pascal }}Controller) Update(c *gin.Context)  {
    id, err := strconv.Atoi(c.Param("id"))

    if err != nil {
      _ = c.Error(core.NewRequestValidationError(err))
      return
    }

    sim := &models.{{ .Name | singular | pascal }}Model{}

    err = c.ShouldBind(sim)
    if err != nil {
      _ = c.Error(core.NewRequestValidationError(err))
      return
    }

    err = controller.app.{{ .Name | singular | pascal}}Repository.Update(c.Request.Context(), id, sim)
    if err != nil {
      _ = c.Error(err)
      return
    }

    item, err := controller.app.{{ .Name | singular | pascal}}Repository.Get(c.Request.Context(), id)
    if err != nil {
      _ = c.Error(err)
      return
    }

    c.JSON(http.StatusOK, item)
  }

  {{ template "go-rest-blocks-swag-get" . }}
  func (controller *{{ .Name | singular | pascal }}Controller) GetPaged(c *gin.Context)  {
    q := models.PageOptions{}

    err := c.ShouldBindQuery(&q)
    if err != nil {
      _ = c.Error(core.NewRequestValidationError(err))
      return
    }

    defPage := 1
    defPageSize := 50

    if q.Page == nil {
      q.Page = &defPage
    }

    if q.PageSize == nil {
      q.PageSize = &defPageSize
    }
    items, err := controller.app.{{ .Name | singular | pascal}}Repository.GetPaged(c.Request.Context(), &q)
    if err != nil {
      _ = c.Error(err)
      return
    }

    c.JSON(http.StatusOK, items)
  }

  {{ range .Children -}}        
  {{ template "go-rest-blocks-swag-get-child" . }}
  func (controller *{{ .OwnerName | singular | pascal }}Controller) Get{{ .NameWithoutPrefix | pascal | plural }}(c *gin.Context)  {
    id, err := strconv.Atoi(c.Param("id"))

    if err != nil {
      c.JSON(http.StatusBadRequest, err)
      return
    }

    items, err := controller.app.{{ .OwnerName | singular | pascal}}Repository.Get{{ .NameWithoutPrefix | pascal | plural }}(c.Request.Context(), id)

    if err != nil {
      _ = c.Error(err)
      return
    }
    
    c.JSON(http.StatusOK, items)
  }
    
  {{ end -}}