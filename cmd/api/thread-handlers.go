package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (app *application) getOneThread(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid thread id parameter"))
		app.errorJSON(w, err)
		return
	}

	thread, err := app.models.DB.Get(id)

	/* thread := models.Thread{
		ID:         id,
		Title:      "Thread 1",
		Content:    "This is the content for thread 1",
		AuthorID:   1,
		AuthorName: "Author 1",
		Upvotes:    0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsSolved:   false,
	} */

	app.writeJSON(w, http.StatusOK, thread, "thread")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := app.models.DB.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, threads, "threads")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.models.DB.CategoriesAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, categories, "categories")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) getAllThreadsByCategory(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	categoryID, err := strconv.Atoi(params.ByName("category_id"))
	if err != nil {
		app.logger.Print(errors.New("invalid category id parameter"))
		app.errorJSON(w, err)
		return
	}

	threads, err := app.models.DB.All(categoryID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, threads, "threads")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) yourThreads(w http.ResponseWriter, r *http.Request) {
    params := httprouter.ParamsFromContext(r.Context())
    authorID, err := strconv.Atoi(params.ByName("author_id"))
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    threads, err := app.models.DB.YourThreads(authorID)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, threads, "threads")
    if err != nil {
        app.errorJSON(w, err)
        return
    }
}

type ThreadPayload struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
	Category   int    `json:"category"`
}

type ReplyPayload struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	AuthorID   int       `json:"author_id"`
	AuthorName string    `json:"author_name"`
	ThreadID   int       `json:"thread_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (app *application) editThread(w http.ResponseWriter, r *http.Request) {
	var payload ThreadPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var thread models.Thread

	if payload.ID != "0" {
		id, err := strconv.Atoi(payload.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
		t, err := app.models.DB.Get(id)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
		thread = *t
		thread.UpdatedAt = time.Now()

	}

	thread.ID, err = strconv.Atoi(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	thread.Title = payload.Title
	thread.Content = payload.Content
	thread.CategoryID = payload.Category
	thread.UpdatedAt = time.Now()

	err = app.models.DB.UpdateThread(thread)
	if err != nil {
		app.errorJSON(w, err)
		return

	}

	ok := jsonResponse{OK: true}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) newThread(w http.ResponseWriter, r *http.Request) {
	var payload ThreadPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var thread models.Thread

	thread.Title = payload.Title
	thread.Content = payload.Content
	thread.AuthorID = payload.AuthorID
	thread.AuthorName = payload.AuthorName
	thread.CategoryID = payload.Category
	thread.UpdatedAt = time.Now()
	thread.CreatedAt = time.Now()

	err = app.models.DB.InsertThread(thread)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ok := jsonResponse{OK: true}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) deleteThread(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.DeleteThread(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ok := jsonResponse{OK: true}

	err = app.writeJSON(w, http.StatusOK, ok, "response")

}

func (app *application) toggleSolved(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid thread id parameter"))
		app.errorJSON(w, err)
		return
	}

	thread, err := app.models.DB.ToggleSolved(id)
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, thread, "thread")
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "Unable to write JSON response", http.StatusInternalServerError)
	}
}

func (app *application) toggleAnswer(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("reply_id"))
	if err != nil {
		app.logger.Print(errors.New("invalid reply id parameter"))
		app.errorJSON(w, err)
		return
	}

	reply, err := app.models.DB.ToggleAnswer(id)
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, reply, "reply")
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "Unable to write JSON response", http.StatusInternalServerError)
	}
}

func (app *application) getReplies(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	threadID, err := strconv.Atoi(params.ByName("thread_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	replies, err := app.models.DB.GetReplies(threadID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, replies, "replies")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) newReply(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	threadID, err := strconv.Atoi(params.ByName("thread_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload ReplyPayload

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var reply models.Reply

	reply.Content = payload.Content
	reply.AuthorID = payload.AuthorID
	reply.AuthorName = payload.AuthorName
	reply.ThreadID = threadID
	reply.CreatedAt = time.Now()

	err = app.models.DB.InsertReply(reply)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ok := jsonResponse{OK: true}
	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) deleteReply(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.DeleteReply(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ok := jsonResponse{OK: true}

	err = app.writeJSON(w, http.StatusOK, ok, "response")

}

func (app *application) starThread(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	userID, err := strconv.Atoi(params.ByName("user_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	threadID, err := strconv.Atoi(params.ByName("thread_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.StarThread(userID, threadID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ok := jsonResponse{OK: true}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) unstarThread(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	userID, err := strconv.Atoi(params.ByName("user_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	threadID, err := strconv.Atoi(params.ByName("thread_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.UnstarThread(userID, threadID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ok := jsonResponse{OK: true}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getStarredThreads (w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userID, err := strconv.Atoi(params.ByName("user_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	threads, err := app.models.DB.GetStarredThreads(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, threads, "threads")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
}