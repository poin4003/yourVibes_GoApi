package implement

//
//import (
//	"context"
//	"github.com/poin4003/yourVibes_GoApi/internal/database"
//	"github.com/poin4003/yourVibes_GoApi/internal/model"
//)
//
//type sUserLogin struct {
//	repo
//}
//
//func NewUserLoginImplement(r *database.Queries) *sUserLogin {
//	return &sUserLogin{r: r}
//}
//
//func (s *sUserLogin) Login(ctx context.Context, in *model.LoginInput) error {
//	return nil
//}
//
//func (s *sUserLogin) Register(ctx context.Context, in *model.RegisterInput) (codeResult int, err error) {
//	//// 1. hash Email
//	//fmt.Printf("Email: %s\n", in.Email)
//	//hashKey := crypto.GetHash(strings.ToLower(in.Email))
//	//
//	//// 2. check user exists in users table
//	//userFound, err := s.r.CheckUserExists(ctx, in.Email)
//	//if err != nil {
//	//	return response.ErrCodeUserHasExists, err
//	//}
//	//
//	//if userFound > 0 {
//	//	return response.ErrCodeUserHasExists, fmt.Errorf("user %s already exists", in.Email)
//	//}
//	//
//	//// 3. Create OTP
//	//userKey := utils.GetUserKey(hashKey)
//	//otpFound, err := global.Rdb.Get(ctx, userKey).Result()
//	//
//	//switch {
//	//case err == redis.Nil:
//	//	fmt.Println("Key does not exist")
//	//case err != nil:
//	//	fmt.Println("Get failed::", err)
//	//	return response.ErrInvalidOTP, err
//	//case otpFound != "":
//	//	return response.ErrCodeOtpNotExists, fmt.Errorf("otp %s already exists but not registered", otpFound)
//	//}
//	//
//	//// 4. Generate OTP
//	//otpNew := random.GenerateSixDigitOtp()
//	//
//	//// 5. save OTP into Redis with expiration time
//	//err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Second).Err()
//	//if err != nil {
//	//	return response.ErrInvalidOTP, err
//	//}
//
//	return 1, nil
//}
