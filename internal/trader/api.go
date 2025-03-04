package trader

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func NewRouter(s *service) *gin.Engine {
	router := gin.Default()

	router.POST("/trader/get-quote", s.HandleGetQuote)
	router.PUT("/trader/accept-quote", s.HandlerAcceptQuote)

	return router
}

// handler logic, can be moved outside trader package

const (
	ErrorCodeInvalidGetQuoteRequest    = "ErrorCodeInvalidGetQuoteRequest"
	ErrorCodeInvalidAcceptQuoteRequest = "ErrorCodeInvalidAcceptQuoteRequest"
	ErrorCodeInvalidQuote              = "ErrorCodeInvalidQuote"
	ErrorCodeNoAvailablePrice          = "ErrorCodeNoAvailablePrice"
	ErrorCodeFailedToSaveQuote         = "ErrorCodeFailedToSaveQuote"
	ErrorCodeFailedToGetQuote          = "ErrorCodeFailedToGetQuote"
)

type errResp struct {
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
	ErrorMsg   string `json:"error_msg"`
}

type getQuoteReq struct {
	UserID   string  `json:"user_id"  binding:"required"`
	Ticker   string  `json:"ticker"   binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Side     string  `json:"side"     binding:"required"`
}

type getQuoteResp struct {
	Quote *Quote `json:"quote"`
}

func (s *service) HandleGetQuote(ctx *gin.Context) {
	var req getQuoteReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  ErrorCodeInvalidGetQuoteRequest,
			ErrorMsg:   err.Error(),
		})

		logrus.Error(err)
		return
	}

	priceStr, err := s.memory.GetPrice(ctx, req.Ticker)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResp{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  ErrorCodeNoAvailablePrice,
			ErrorMsg:   err.Error(),
		})

		logrus.Error(err)
		return
	}

	price, _ := decimal.NewFromString(priceStr)
	quantity := decimal.NewFromFloat(req.Quantity)

	quote := Quote{
		UserID:    uuid.NewString(),
		QuoteID:   uuid.NewString(),
		Ticker:    req.Ticker,
		Quantity:  quantity,
		Price:     price,
		Amount:    price.Mul(quantity),
		Side:      req.Side,
		CreatedAt: time.Now(),
	}

	if err := s.memory.SaveQuote(ctx, quote); err != nil {
		logrus.Error(err)

		ctx.JSON(http.StatusInternalServerError, errResp{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  ErrorCodeFailedToSaveQuote,
			ErrorMsg:   err.Error(),
		})

		logrus.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, getQuoteResp{
		Quote: &quote,
	})
}

type accepQuoteReq struct {
	UserID  string `json:"user_id"  binding:"required"`
	QuoteID string `json:"quote_id" binding:"required"`
}

type accepQuoteResp struct {
	UserID  string `json:"user_id"`
	QuoteID string `json:"quote_id"`
	Msg     string `json:"msg"`
}

func (s *service) HandlerAcceptQuote(ctx *gin.Context) {
	var req accepQuoteReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  ErrorCodeInvalidAcceptQuoteRequest,
			ErrorMsg:   err.Error(),
		})

		logrus.Error(err)
		return
	}

	q, err := s.memory.GetQuote(ctx, req.UserID, req.QuoteID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResp{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  ErrorCodeFailedToGetQuote,
			ErrorMsg:   err.Error(),
		})

		logrus.Error(err)
		return
	}

	if valid, err := q.IsValid(); !valid || err != nil {
		ctx.JSON(http.StatusInternalServerError, errResp{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  "invalid quote or expired",
			ErrorMsg:   err.Error(),
		})

		logrus.Error(err)
		return
	}

	// save quote
	// save trade

	ctx.JSON(http.StatusOK, accepQuoteResp{
		QuoteID: q.QuoteID,
		UserID:  q.UserID,
		Msg:     "quote has been successfully accepted",
	})
}
