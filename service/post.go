package service

import (
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"go-blog/common"
	"go-blog/conf"
	"go-blog/entity"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"time"
	"xorm.io/xorm"
)

func ConsolePostCount(limit int, offset int, isTrash bool) (count int64, err error) {
	post := new(entity.ZPosts)
	if isTrash {
		count, err = conf.SqlServer.Unscoped().Where("`deleted_at` IS NOT NULL OR `deleted_at`=?", "0001-01-01 00:00:00").Desc("id").Limit(limit, offset).Count(post)
	} else {
		count, err = conf.SqlServer.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Desc("id").Limit(limit, offset).Count(post)
	}
	if err != nil {

		return 0, err
	}
	return count, nil
}

func ConsolePostIndex(limit int, offset int, isTrash bool) (postListArr []*common.ConsolePostList, err error) {
	post := new(entity.ZPosts)
	var rows *xorm.Rows
	if isTrash {
		rows, err = conf.SqlServer.Unscoped().Where("`deleted_at` IS NOT NULL OR `deleted_at`=?", "0001-01-01 00:00:00").Desc("id").Limit(limit, offset).Rows(post)
	} else {
		rows, err = conf.SqlServer.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Desc("id").Limit(limit, offset).Rows(post)
	}

	if err != nil {

		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		//post
		post := new(entity.ZPosts)
		err = rows.Scan(post)
		if err != nil {

			return nil, err
		}

		consolePost := common.ConsolePost{
			Id:        post.Id,
			Uid:       post.Uid,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		//category
		cates, err := GetPostCateByPostId(post.Id)
		if err != nil {

			return nil, err
		}
		consoleCate := common.ConsoleCate{
			Id:          cates.Id,
			Name:        cates.Name,
			DisplayName: cates.DisplayName,
			SeoDesc:     cates.SeoDesc,
		}

		//tag
		tagIds, err := GetPostTagsByPostId(post.Id)
		if err != nil {

			return nil, err
		}
		tags, err := GetTagsByIds(tagIds)
		if err != nil {

			return nil, err
		}
		var consoleTags []common.ConsoleTag
		for _, v := range tags {
			consoleTag := common.ConsoleTag{
				Id:          v.Id,
				Name:        v.Name,
				DisplayName: v.DisplayName,
				SeoDesc:     v.SeoDesc,
				Num:         v.Num,
			}
			consoleTags = append(consoleTags, consoleTag)
		}

		//view
		view, err := PostView(post.Id)
		if err != nil {

			return nil, err
		}
		consoleView := common.ConsoleView{
			Num: view.Num,
		}

		//user
		user, err := GetUserById(post.UserId)
		if err != nil {

			return nil, err
		}
		consoleUser := common.ConsoleUser{
			Id:     user.Id,
			Name:   user.Name,
			Email:  user.Email,
			Status: user.Status,
		}

		postList := common.ConsolePostList{
			Post:     consolePost,
			Category: consoleCate,
			Tags:     consoleTags,
			TagsLen:  len(tags) - 1,
			View:     consoleView,
			Author:   consoleUser,
		}
		postListArr = append(postListArr, &postList)
	}

	return postListArr, nil
}

func PostView(postId int) (*entity.ZPostViews, error) {
	postV := new(entity.ZPostViews)
	_, err := conf.SqlServer.Where("post_id = ?", postId).Cols("num").Get(postV)
	if err != nil {

	}
	return postV, nil
}

func PostStore(ps common.PostStore, userId int) {
	postCreate := &entity.ZPosts{
		Title:    ps.Title,
		UserId:   userId,
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	postCreate.Content = string(html)

	session := conf.SqlServer.NewSession()
	defer session.Close()
	affected, err := session.Insert(postCreate)
	if err != nil {

		_ = session.Rollback()
		return
	}
	if affected < 1 {

		_ = session.Rollback()
		return
	}

	if ps.Category > 0 {
		postCateCreate := entity.ZPostCate{
			PostId: postCreate.Id,
			CateId: ps.Category,
		}
		affected, err := session.Insert(postCateCreate)
		if err != nil {

			_ = session.Rollback()
			return
		}

		if affected < 1 {

			_ = session.Rollback()
			return
		}
	}

	if len(ps.Tags) > 0 {
		for _, v := range ps.Tags {
			postTagCreate := entity.ZPostTag{
				PostId: postCreate.Id,
				TagId:  v,
			}
			affected, err := session.Insert(postTagCreate)
			if err != nil {

				_ = session.Rollback()
				return
			}
			if affected < 1 {

				_ = session.Rollback()
				return
			}

			affected, err = session.ID(v).Incr("num").Update(entity.ZTags{})
			if err != nil {

				_ = session.Rollback()
				return
			}
			if affected < 1 {

				_ = session.Rollback()
				return
			}
		}
	}

	postView := entity.ZPostViews{
		PostId: postCreate.Id,
		Num:    1,
	}

	affected, err = session.Insert(postView)
	if err != nil {

		_ = session.Rollback()
		return
	}

	if affected < 1 {

		_ = session.Rollback()
		return
	}

	_ = session.Commit()

	uid, err := conf.ZHashId.Encode([]int{postCreate.Id})
	if err != nil {

		return
	}

	newPostCreate := entity.ZPosts{
		Uid: uid,
	}
	affected, err = session.Where("id = ?", postCreate.Id).Update(newPostCreate)
	if err != nil {

		return
	}

	if affected < 1 {

		return
	}

	return
}

func PostDetail(postId int) (p *entity.ZPosts, err error) {
	post := new(entity.ZPosts)
	_, err = conf.SqlServer.Where("id = ?", postId).Get(post)
	if err != nil {

		return post, err
	}
	return post, nil
}

func IndexPostDetailDao(postId int) (postDetail common.IndexPostDetail, err error) {
	post := new(entity.ZPosts)
	_, err = conf.SqlServer.Where("id = ?", postId).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Get(post)
	if err != nil {

		return
	}
	if post.Id <= 0 {
		return postDetail, errors.New("Post do not exists ")
	}
	Post := common.IndexPost{
		Id:        post.Id,
		Uid:       post.Uid,
		Title:     post.Title,
		Summary:   post.Summary,
		Original:  post.Original,
		Content:   template.HTML(post.Content),
		Password:  post.Password,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	tags, err := PostIdTags(postId)
	if err != nil {

		return
	}
	var Tags []common.ConsoleTag
	for _, v := range tags {
		consoleTag := common.ConsoleTag{
			Id:          v.Id,
			Name:        v.Name,
			DisplayName: v.DisplayName,
			SeoDesc:     v.SeoDesc,
			Num:         v.Num,
		}
		Tags = append(Tags, consoleTag)
	}

	cate, err := PostCates(postId)
	if err != nil {

		return
	}
	Cate := common.ConsoleCate{
		Id:          cate.Id,
		Name:        cate.Name,
		DisplayName: cate.DisplayName,
		SeoDesc:     cate.SeoDesc,
	}

	//view
	view, err := PostView(post.Id)
	if err != nil {

		return
	}
	View := common.ConsoleView{
		Num: view.Num,
	}

	//user
	user, err := GetUserById(post.UserId)
	if err != nil {

		return
	}
	Author := common.ConsoleUser{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}

	// last post
	lastPost, err := LastPost(postId)
	if err != nil {

		return
	}

	// next post
	nextPost, err := NextPost(postId)
	if err != nil {

		return
	}

	postDetail = common.IndexPostDetail{
		Post:     Post,
		Category: Cate,
		Tags:     Tags,
		TagsLen:  len(Tags) - 1,
		View:     View,
		Author:   Author,
		LastPost: lastPost,
		NextPost: nextPost,
	}

	return postDetail, nil
}

func LastPost(postId int) (post *entity.ZPosts, err error) {
	post = new(entity.ZPosts)
	_, err = conf.SqlServer.Where("id < ?", postId).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Desc("id").Get(post)
	return
}

func NextPost(postId int) (post *entity.ZPosts, err error) {
	post = new(entity.ZPosts)
	_, err = conf.SqlServer.Where("id > ?", postId).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Asc("id").Get(post)
	return
}

func PostIdTags(postId int) (tags []*entity.ZTags, err error) {
	tagIds, err := PostIdTag(postId)
	if err != nil {

		return
	}
	//tags = make([]entity.ZTags,0)
	err = conf.SqlServer.In("id", tagIds).Find(&tags)
	return
}

func PostIdTag(postId int) (tagIds []int, err error) {
	postTag := make([]entity.ZPostTag, 0)
	err = conf.SqlServer.Where("post_id = ?", postId).Find(&postTag)
	if err != nil {

		return
	}

	for _, v := range postTag {
		tagIds = append(tagIds, v.TagId)
	}
	return tagIds, nil
}

func PostCates(postId int) (cate *entity.ZCategories, err error) {
	cateId, err := PostCate(postId)
	if err != nil {

		return
	}
	cate = new(entity.ZCategories)
	_, err = conf.SqlServer.ID(cateId).Get(cate)
	if err != nil {

		return
	}
	return
}

func PostCate(postId int) (int, error) {
	postCate := new(entity.ZPostCate)
	_, err := conf.SqlServer.Where("post_id = ?", postId).Get(postCate)
	if err != nil {

		return 0, err
	}
	return postCate.CateId, nil
}

func PostUpdate(postId int, ps common.PostStore) {
	postUpdate := &entity.ZPosts{
		Title:    ps.Title,
		UserId:   1,
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	postUpdate.Content = string(html)

	session := conf.SqlServer.NewSession()
	defer session.Close()
	affected, err := session.Where("id = ?", postId).Update(postUpdate)
	if err != nil {

		_ = session.Rollback()
		return
	}
	if affected < 1 {

		_ = session.Rollback()
		return
	}

	postCate := new(entity.ZPostCate)
	_, err = session.Where("post_id = ?", postId).Delete(postCate)
	if err != nil {

		_ = session.Rollback()
		return
	}

	if ps.Category > 0 {
		postCateCreate := entity.ZPostCate{
			PostId: postId,
			CateId: ps.Category,
		}

		affected, err := session.Insert(postCateCreate)
		if err != nil {

			_ = session.Rollback()
			return
		}

		if affected < 1 {

			_ = session.Rollback()
			return
		}
	}

	postTag := make([]entity.ZPostTag, 0)
	err = session.Where("post_id = ?", postId).Find(&postTag)

	if err != nil {

		_ = session.Rollback()
		return
	}

	if len(postTag) > 0 {
		for _, v := range postTag {
			affected, err = session.ID(v.TagId).Decr("num").Update(entity.ZTags{})
			if err != nil {

				_ = session.Rollback()
				return
			}
			if affected < 1 {

				_ = session.Rollback()
				return
			}
		}

		_, err = session.Where("post_id = ?", postId).Delete(new(entity.ZPostTag))

		if err != nil {

			_ = session.Rollback()
			return
		}
	}

	if len(ps.Tags) > 0 {
		for _, v := range ps.Tags {
			postTagCreate := entity.ZPostTag{
				PostId: postId,
				TagId:  v,
			}
			affected, err := session.Insert(postTagCreate)
			if err != nil {

				_ = session.Rollback()
				return
			}
			if affected < 1 {

				_ = session.Rollback()
				return
			}
			affected, err = session.ID(v).Incr("num").Update(entity.ZTags{})
			if err != nil {

				_ = session.Rollback()
				return
			}
			if affected < 1 {

				_ = session.Rollback()
				return
			}
		}
	}
	_ = session.Commit()

	return
}

func PostDestroy(postId int) (bool, error) {
	post := new(entity.ZPosts)
	toBeCharge := time.Now().Format(time.RFC3339)
	timeLayout := time.RFC3339
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, toBeCharge, loc)
	post.DeletedAt = &theTime
	_, err = conf.SqlServer.ID(postId).Update(post)
	if err != nil {

		return false, err
	}
	return true, nil
}

func PostUnTrash(postId int) (bool, error) {
	post := new(entity.ZPosts)
	theTime, _ := time.Parse("2006-01-02 15:04:05", "")
	post.DeletedAt = &theTime
	_, err := conf.SqlServer.ID(postId).Update(post)
	if err != nil {

		return false, err
	}
	return true, nil
}

func PostCnt() (cnt int64, err error) {
	post := new(entity.ZPosts)
	cnt, err = conf.SqlServer.Count(post)
	return
}

func PostTagListCount(tagId int, limit int, offset int) (count int64, err error) {
	postTag := new(entity.ZPostTag)
	count, err = conf.SqlServer.Where("tag_id = ?", tagId).Desc("id").Limit(limit, offset).Count(postTag)
	if err != nil {

		return 0, err
	}
	return
}

func PostTagList(tagId int, limit int, offset int) (postListArr []*common.ConsolePostList, err error) {
	postTag := new(entity.ZPostTag)
	rows, err := conf.SqlServer.Where("tag_id = ?", tagId).Desc("id").Limit(limit, offset).Rows(postTag)

	if err != nil {

		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		//post
		postTag := new(entity.ZPostTag)
		err = rows.Scan(postTag)
		if err != nil {

			return nil, err
		}

		post := new(entity.ZPosts)
		_, err = conf.SqlServer.ID(postTag.PostId).Get(post)

		consolePost := common.ConsolePost{
			Id:        post.Id,
			Uid:       post.Uid,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		postList := common.ConsolePostList{
			Post: consolePost,
		}
		postListArr = append(postListArr, &postList)
	}

	return postListArr, nil
}

func PostCateListCount(cateId int, limit int, offset int) (count int64, err error) {
	postCate := new(entity.ZPostCate)
	count, err = conf.SqlServer.Where("cate_id = ?", cateId).Desc("id").Limit(limit, offset).Count(postCate)
	if err != nil {

		return 0, err
	}
	return
}

func PostCateList(cateId int, limit int, offset int) (postListArr []*common.ConsolePostList, err error) {
	postCate := new(entity.ZPostCate)
	rows, err := conf.SqlServer.Where("cate_id = ?", cateId).Desc("id").Limit(limit, offset).Rows(postCate)

	if err != nil {

		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		//post
		postCate := new(entity.ZPostCate)
		err = rows.Scan(postCate)
		if err != nil {

			return nil, err
		}

		post := new(entity.ZPosts)
		_, err = conf.SqlServer.ID(postCate.PostId).Get(post)

		consolePost := common.ConsolePost{
			Id:        post.Id,
			Uid:       post.Uid,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		postList := common.ConsolePostList{
			Post: consolePost,
		}
		postListArr = append(postListArr, &postList)
	}

	return postListArr, nil
}
