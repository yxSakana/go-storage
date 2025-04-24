package fileDownload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadRangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadRangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadRangeLogic {
	return &DownloadRangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadRangeLogic) DownloadRange(req *types.DownloadRangeReq, w http.ResponseWriter, r *http.Request) (resp *types.DownloadRangeResp, err error) {
	fileMetaRet, err := l.svcCtx.FileMetaModel.FindOne(l.ctx, req.FileId)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fileMetaRet.Path)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+url.QueryEscape(fi.Name()))

	fileSize := fi.Size()
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
		_, err = io.Copy(w, file)
		return resp, err
	}

	// 解析 Range 头
	start, end, err := parseRange(rangeHeader, fileSize)
	if err != nil {
		return nil, fmt.Errorf("invalid range")
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.WriteHeader(http.StatusPartialContent) // 206

	// 读取指定范围
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(w, file, end-start+1)
	return resp, err
}

func parseRange(header string, fileSize int64) (start, end int64, err error) {
	if !strings.HasPrefix(header, "bytes=") {
		return 0, 0, fmt.Errorf("invalid range")
	}

	ranges := strings.TrimPrefix(header, "bytes=")
	parts := strings.Split(ranges, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range parts")
	}

	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid range start")
	}

	if parts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid range end")
		}
	}

	if start > end || end >= fileSize {
		return 0, 0, fmt.Errorf("invalid range position")
	}

	return start, end, nil
}
