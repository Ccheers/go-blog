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
	"time"

	"github.com/snabb/sitemap"
)

// GetSiteMap 获取站点sitemap
func GetSiteMap() (sm *sitemap.Sitemap) {
	sm = sitemap.New()

	staticSitemap(sm)
	postSiteMap(sm)

	return sm
}

func staticSitemap(sm *sitemap.Sitemap) {
	today := time.Now()
	sm.Add(&sitemap.URL{
		Loc:        "https://www.ericcai.fun",
		LastMod:    &today,
		ChangeFreq: sitemap.Daily,
	})
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
		sm.Add(&sitemap.URL{
			Loc:        fmt.Sprintf("https://www.ericcai.fun/detail/%d", post.Id),
			LastMod:    &today,
			ChangeFreq: sitemap.Daily,
		})
	}
}
