package elasticache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/elasticache"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfelasticache "github.com/hashicorp/terraform-provider-aws/internal/service/elasticache"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccElastiCacheUserGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var userGroup elasticache.UserGroup
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserGroupExists(ctx, resourceName, &userGroup),
					resource.TestCheckResourceAttr(resourceName, "user_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_group_id", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccElastiCacheUserGroup_update(t *testing.T) {
	ctx := acctest.Context(t)
	var userGroup elasticache.UserGroup
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserGroupExists(ctx, resourceName, &userGroup),
					resource.TestCheckResourceAttr(resourceName, "user_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_group_id", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
				),
			},
			{
				Config: testAccUserGroupConfig_multiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserGroupExists(ctx, resourceName, &userGroup),
					resource.TestCheckResourceAttr(resourceName, "user_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "user_group_id", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
				),
			},
			{
				Config: testAccUserGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserGroupExists(ctx, resourceName, &userGroup),
					resource.TestCheckResourceAttr(resourceName, "user_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_group_id", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
				),
			},
		},
	})
}

func TestAccElastiCacheUserGroup_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var userGroup elasticache.UserGroup
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserGroupExists(ctx, resourceName, &userGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccUserGroupConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserGroupExists(ctx, resourceName, &userGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccUserGroupConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserGroupExists(ctx, resourceName, &userGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccElastiCacheUserGroup_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var userGroup elasticache.UserGroup
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserGroupExists(ctx, resourceName, &userGroup),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfelasticache.ResourceUserGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckUserGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ElastiCacheConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_elasticache_user_group" {
				continue
			}

			_, err := tfelasticache.FindUserGroupByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("ElastiCache User Group (%s) still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckUserGroupExists(ctx context.Context, n string, v *elasticache.UserGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ElastiCache User Group ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ElastiCacheConn()

		output, err := tfelasticache.FindUserGroupByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccUserGroupConfig_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test1" {
  user_id       = "%[1]s-1"
  user_name     = "default"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"
  passwords     = ["password123456789"]
}

resource "aws_elasticache_user" "test2" {
  user_id       = "%[1]s-2"
  user_name     = "username1"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"
  passwords     = ["password123456789"]
}

resource "aws_elasticache_user_group" "test" {
  user_group_id = %[1]q
  engine        = "REDIS"
  user_ids      = [aws_elasticache_user.test1.user_id]
}
`, rName))
}

func testAccUserGroupConfig_multiple(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test1" {
  user_id       = "%[1]s-1"
  user_name     = "default"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"
  passwords     = ["password123456789"]
}

resource "aws_elasticache_user" "test2" {
  user_id       = "%[1]s-2"
  user_name     = "username1"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"
  passwords     = ["password123456789"]
}

resource "aws_elasticache_user_group" "test" {
  user_group_id = %[1]q
  engine        = "REDIS"
  user_ids      = [aws_elasticache_user.test1.user_id, aws_elasticache_user.test2.user_id]
}
`, rName))
}

func testAccUserGroupConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test1" {
  user_id       = "%[1]s-1"
  user_name     = "default"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"
  passwords     = ["password123456789"]
}
resource "aws_elasticache_user_group" "test" {
  user_group_id = %[1]q
  engine        = "REDIS"
  user_ids      = [aws_elasticache_user.test1.user_id]
  tags = {
    %[2]s = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccUserGroupConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test1" {
  user_id       = "%[1]s-1"
  user_name     = "default"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"
  passwords     = ["password123456789"]
}
resource "aws_elasticache_user_group" "test" {
  user_group_id = %[1]q
  engine        = "REDIS"
  user_ids      = [aws_elasticache_user.test1.user_id]
  tags = {
    %[2]s = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}
