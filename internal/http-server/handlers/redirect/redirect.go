// Мой код!!!!!!!!!!!!

// package redirect

// import (
// 	"errors"
// 	"log/slog"
// 	"net/http"
// 	"url-shortener/internal/lib/logger/sl"
// 	"url-shortener/internal/storage"

// 	// "url-shortener/internal/http-server/middleware/logger"
// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// )

// type URLGetter interface {
// 	GetURL(alias string) (string, error)
// }

// func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		const op = "handlers.url.redirect.New"

// 		log = log.With(
// 			slog.String("op", op),
// 			slog.String("request_id", middleware.GetReqID(r.Context())),
// 		)

// 		alias := chi.URLParam(r, "alias")

// 		if alias == "" {
// 			log.Error("Alias is required")
// 			http.Error(w, "Alias is required", http.StatusBadRequest)
// 			return
// 		}

// 		log.Info("alias from URL", slog.String("alias", alias))

// 		url, err := urlGetter.GetURL(alias)
// 		if errors.Is(err, storage.ErrURLNotFound) {
// 			log.Info("url not found", slog.String("url", url))

// 			http.Error(w, "URL not found", http.StatusNotFound)

// 			return
// 		}
// 		if err != nil {
// 			log.Error("failed to get url", sl.Err(err))

// 			http.Error(w, "Failed to get URL", http.StatusInternalServerError)

// 			return
// 		}
// 		log.Info("redirecting to full URL", slog.String("url", url))
// 		http.Redirect(w, r, url, http.StatusFound)

//		}
//	}
package redirect

import (
	"errors"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

// URLGetter is an interface for getting url by alias.
//

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		// redirect to found url
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
