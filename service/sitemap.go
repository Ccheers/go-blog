package service

//
//
//
//        ***************************     ***************************         *********      ************************
//      *****************************    ******************************      *********      *************************
//     *****************************     *******************************     *********     *************************
//    *********                         *********                *******    *********     *********
//    ********                          *********               ********    *********     ********
//   ********     ******************   *********  *********************    *********     *********
//   ********     *****************    *********  ********************     *********     ********
//  ********      ****************    *********     ****************      *********     *********
//  ********                          *********      ********             *********     ********
// *********                         *********         ******            *********     *********
// ******************************    *********          *******          *********     *************************
//  ****************************    *********            *******        *********      *************************
//    **************************    *********              ******       *********         *********************
//
//
import (
	"fmt"
	"go-blog/conf"
	"go-blog/entity"
	"go-blog/sitemap"
	"time"
)

// GetSiteMap 获取站点sitemap
func GetSiteMap() (sm *sitemap.Sitemap) {
	sm = sitemap.NewSitemap()

	staticSitemap(sm)
	postSiteMap(sm)

	return sm
}

func staticSitemap(sm *sitemap.Sitemap) {
	today := time.Now()
	sm.Urls = append(sm.Urls, sitemap.NewUrl("https://www.ericcai.fun", today, sitemap.ChangefreqDaily, 1))
}

func postSiteMap(sm *sitemap.Sitemap) {
	today := time.Now()

	post := new(entity.ZPosts)
	rows, err := conf.SqlServer.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Rows(post)
	if err != nil {
		return
	}
	for rows.Next() {
		//post
		post := new(entity.ZPosts)
		err = rows.Scan(post)
		if err != nil {
			break
		}
		sm.Urls = append(sm.Urls, sitemap.NewUrl(
			fmt.Sprintf("https://www.ericcai.fun/detail/%d", post.Id),
			today,
			sitemap.ChangefreqDaily,
			0.7,
		))
	}
}
