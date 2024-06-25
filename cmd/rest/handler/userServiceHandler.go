package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ngc5/model"
	"ngc5/pb"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	UserGRPC pb.UserServiceClient
}

func (h *UserHandler) AddUser(e echo.Context) error {
	var payload model.User

	err := e.Bind(&payload)
	if err != nil {
		log.Println("ERR BIND:", err)
		return e.JSON(http.StatusInternalServerError, "Error Binding")
	}

	//validate payload
	if payload.Name == "" {
		return e.JSON(http.StatusBadRequest, "Error or missing param")
	}

	in := pb.AddRequest{
		Name: payload.Name,
	}

	response, err := h.UserGRPC.AddUser(context.TODO(), &in)
	if err != nil {
		log.Println("ERR USER SERVICE: ", err)
		return e.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return e.JSON(http.StatusCreated, response)
}

func (h *UserHandler) GetUser(e echo.Context) error {
	var payload model.User

	err := e.Bind(&payload)
	if err != nil {
		log.Println("ERR BIND:", err)
		return e.JSON(http.StatusInternalServerError, "Error Binding")
	}

	//validate payload
	if payload.Name == "" {
		return e.JSON(http.StatusBadRequest, "Error or missing param")
	}

	in := &pb.GetRequest{
		Name: payload.Name,
	}

	res, err := h.UserGRPC.GetUser(context.TODO(), in)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// Extract the error code and handle accordingly
			switch st.Code() {
			case codes.NotFound:
				return e.JSON(http.StatusNotFound, st.Message())
			default:
				fmt.Printf("Unexpected error: %v\n", st.Message())
			}
		} else {
			fmt.Printf("Non-gRPC error: %v\n", err)
			log.Println("ERR USER SERVICE: ", err)
			return e.JSON(http.StatusInternalServerError, "Internal Server Error")
		}
	}
	return e.JSON(http.StatusOK, res)
}
