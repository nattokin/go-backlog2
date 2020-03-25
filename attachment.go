package backlog

import (
	"encoding/json"
	"strconv"
)

// AttachmentService hs methods for attachment.
type baseAttachmentService struct {
	clientMethod *clientMethod
}

// Uploade uploads a any file to the space.
//
// File's path and name are must not empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/post-attachment-file
func (s *baseAttachmentService) Uploade(fPath, fName string) (*Attachment, error) {
	spath := "space/attachment"
	resp, err := s.clientMethod.Uploade(spath, fPath, fName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Attachment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *baseAttachmentService) list(spath string) ([]*Attachment, error) {
	resp, err := s.clientMethod.Get(spath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Attachment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *baseAttachmentService) remove(spath string) (*Attachment, error) {
	resp, err := s.clientMethod.Delete(spath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Attachment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// WikiAttachmentService hs methods for attachment file of wiki.
type WikiAttachmentService struct {
	*baseAttachmentService
}

// Attach attachs files uploaded to space to the wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/attach-file-to-wiki
func (s *WikiAttachmentService) Attach(wikiID int, attachmentIDs []int) ([]*Attachment, error) {
	params := newRequestParams()
	for _, id := range attachmentIDs {
		params.Add("attachmentId[]", strconv.Itoa(id))
	}
	spath := "wikis/" + strconv.Itoa(wikiID) + "/attachments"
	resp, err := s.clientMethod.Post(spath, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Attachment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

// List returns a list of all attachments in the wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-wiki-attachments
func (s *WikiAttachmentService) List(wikiID int) ([]*Attachment, error) {
	spath := "wikis/" + strconv.Itoa(wikiID) + "/attachments"
	return s.list(spath)
}

// Remove removes a file attached to the wiki.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-wiki-attachment
func (s *WikiAttachmentService) Remove(wikiID, attachmentID int) (*Attachment, error) {
	spath := "wikis/" + strconv.Itoa(wikiID) + "/attachments/" + strconv.Itoa(attachmentID)
	return s.remove(spath)
}

// IssueAttachmentService hs methods for attachment file of issue.
type IssueAttachmentService struct {
	*baseAttachmentService
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(issueIDOrKey string) ([]*Attachment, error) {
	spath := "issues/" + issueIDOrKey + "/attachments"
	return s.list(spath)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(issueIDOrKey string, attachmentID int) (*Attachment, error) {
	spath := "issues/" + issueIDOrKey + "/attachments/" + strconv.Itoa(attachmentID)
	return s.remove(spath)
}

// PullRequestAttachmentService hs methods for attachment file of pull request.
type PullRequestAttachmentService struct {
	*baseAttachmentService
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(projectIDOrKey, repoIDOrName string, prNumber int) ([]*Attachment, error) {
	spath := "projects/" + projectIDOrKey + "/git/repositories/" + repoIDOrName + "/pullRequests/" + strconv.Itoa(prNumber) + "/attachments"
	return s.list(spath)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(projectIDOrKey, repoIDOrName string, prNumber int) (*Attachment, error) {
	spath := "projects/" + projectIDOrKey + "/git/repositories/" + repoIDOrName + "/pullRequests/" + strconv.Itoa(prNumber) + "/attachments"
	return s.remove(spath)
}
