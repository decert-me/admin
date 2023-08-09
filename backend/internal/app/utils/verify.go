package utils

var (
	IDVerify = Rules{"ID": {NotEmpty()}}

	PageInfoVerify = Rules{"Page": {Gt("0")}, "PageSize": {NotEmpty(), Le("100")}}
	// 用户
	LoginVerify           = Rules{"CaptchaId": {NotEmpty()}, "Captcha": {NotEmpty()}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}
	LoginVerifyNotCaptcha = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify        = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}}
	ChangePasswordVerify  = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}, "RepeatNewPassword": {NotEmpty()}}
	AuthorityVerify       = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}, "ParentId": {NotEmpty()}}
	AuthorityIdVerify     = Rules{"AuthorityId": {NotEmpty()}}
	AuthorityDataVerify   = Rules{"AuthorityId": {NotEmpty()}, "AuthoritySourceId": {NotEmpty()}}

	// 消息
	CreateNewsVerify = Rules{"Title": {NotEmpty()}, "Content": {NotEmpty()}}
	ModifyNewsVerify = Rules{"id": {NotEmpty()}, "Title": {NotEmpty()}, "Content": {NotEmpty()}}

	// 产品评论
	ModifyProductCommentaryVerify = Rules{"id": {NotEmpty()}, "Title": {NotEmpty()}, "Content": {NotEmpty()}}

	// 研究报告
	ModifyResearchReportVerify = Rules{"id": {NotEmpty()}, "Title": {NotEmpty()}, "Content": {NotEmpty()}}

	// 特别报告
	// SpecialEditionsPermissionVerify = Rules{"Permission": {Eq("2")}} //TODO: 暂时放宽
	SpecialEditionsPermissionVerify = Rules{}
	ModifySpecialEditionsVerify     = Rules{"id": {NotEmpty()}, "Title": {NotEmpty()}, "Content": {NotEmpty()}}

	// banner
	CreateBannerVerify = Rules{"Link": {NotEmpty()}}
	DeleteBannerVerify = Rules{"ID": {Gt("1")}}

	// insight
	InsightMainVerify = Rules{"Image": {NotEmpty()}, "File": {NotEmpty()}}
	InsightVerify     = Rules{"Title": {NotEmpty()}, "Language": {NotEmpty()}}

	// production
	CreateProductionVerify = Rules{"Image": {NotEmpty()}}
	DeleteProductionVerify = Rules{"ID": {Gt("0")}}

	// AMC
	CreateAMCVerify = Rules{"Tiele": {NotEmpty()}, "Description": {NotEmpty()}}
	DeleteAMCVerify = Rules{"ID": {Gt("0")}}
)
