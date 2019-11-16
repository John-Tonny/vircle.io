package model

import (
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego/orm"
	"vircle.io/mnhost/model"

	//"vircle.io/mnhost/utils"

	pb "vircle.io/mnhost/interface/out/user"
)

type Repository interface {
	Get(id string) (*pb.User, error)
	GetAll() ([]*pb.User, error)
	Create(*pb.User) error
	GetByMobile(mobile string) (*pb.User, error)
	Close()
}

type UserRepository struct {
	session *mgo.Session
}

const (
	DB_NAME        = "MicroServicePractice"
	CON_COLLECTION = "users"
)

func GetUserRepository(session *mgo.Session) *UserRepository {
	return &UserRepository{session: session}
}

func (repo *UserRepository) Get(id string) (*pb.User, error) {
	iid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("ID format is wrong")
	}

	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err = qs.Filter("id", iid).One(&user)
	if err != nil {
		return nil, err
	}

	pbUser := User2PBUser(&user)
	return &pbUser, nil
}

func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	nums, err := qs.All(&users)
	if err != nil {
		return nil, err
	}
	if nums == 0 {
		return nil, errors.New("no data")
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUser := User2PBUser(&user)
		pbUsers[i] = &pbUser
	}
	return pbUsers, nil
}

func (repo *UserRepository) Create(u *pb.User) error {
	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("mobile", u.Mobile).One(&user)
	if err == nil {
		return errors.New("exist")
	}

	o = orm.NewOrm()
	fmt.Println(u.Name)
	fmt.Println(u.Password)
	fmt.Println(u.Mobile)

	user = PBUser2User(u)
	fmt.Println(user.Name)
	fmt.Println(user.Password_hash)
	fmt.Println(user.Mobile)

	userId, err := o.Insert(&user)
	if err != nil {
		return err
	}

	u.Id = strconv.FormatInt(userId, 10)
	fmt.Println("userId:", u.Id)

	return nil
}

func (repo *UserRepository) GetByMobile(mobile string) (*pb.User, error) {
	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("mobile", mobile).One(&user)
	if err != nil {
		return nil, err
	}

	pbUser := User2PBUser(&user)
	return &pbUser, nil
}

// 关闭连接
func (repo *UserRepository) Close() {
	// Close() 会在每次查询结束的时候关闭会话
	// Mgo 会在启动的时候生成一个 "主" 会话
	// 你可以使用 Copy() 直接从主会话复制出新会话来执行，即每个查询都会有自己的数据库会话
	// 同时每个会话都有自己连接到数据库的 socket 及错误处理，这么做既安全又高效
	// 如果只使用一个连接到数据库的主 socket 来执行查询，那很多请求处理都会阻塞
	// Mgo 因此能在不使用锁的情况下完美处理并发请求
	// 不过弊端就是，每次查询结束之后，必须确保数据库会话要手动 Close
	// 否则将建立过多无用的连接，白白浪费数据库资源
	repo.session.Close()
}

// 返回所有货物信息
func (repo *UserRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(CON_COLLECTION)
}

func PBUser2User(u *pb.User) models.User {
	data := models.User{
		Name:          u.Name,
		Password_hash: u.Password,
		Mobile:        u.Mobile,
		//Real_name:     u.RealName,
		//Id_card:       u.IdCard,
	}
	return data
}

func User2PBUser(u *models.User) pb.User {
	return pb.User{
		Id:       strconv.Itoa(u.Id),
		Name:     u.Name,
		Password: u.Password_hash,
		Mobile:   u.Mobile,
		RealName: u.Real_name,
		IdCard:   u.Id_card,
	}
}
