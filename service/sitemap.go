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
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"go-blog/conf"
	"go-blog/entity"
	"time"
)

//获取站点sitemap
func GetSiteMap() (sm *stm.Sitemap) {
	sm = stm.NewSitemap(1)

	sm.SetDefaultHost("https://www.ericcai.fun")

	staticSitemap(sm)
	postSiteMap(sm)

	return sm
}

func staticSitemap(sm *stm.Sitemap) {
	today := time.Now().Format("2006-01-02")
	sm.Add(
		stm.URL{
			{"loc", "/"},
			{"changefreq", "daily"},
			{"priority", 0.7},
			{"lastmod", today},
		},
	)
}

func postSiteMap(sm *stm.Sitemap) {
	today := time.Now().Format("2006-01-02")

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
		sm.Add(
			stm.URL{
				{"loc", fmt.Sprintf("/detail/%d", post.Id)},
				{"changefreq", "daily"},
				{"priority", 0.7},
				{"lastmod", today},
			},
		)
	}
}
