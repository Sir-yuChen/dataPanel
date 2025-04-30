export type UserLogin = {
    loginType: string;
    loginPassword?: string;
    mobileNo?: string;
    captcha?: string;
    iv?: string;
    encryptedData?: string;
}

export type LoginReturnData = {
    token: string,
    user: UserInfo
}
export type UserInfo = {
    accountType: string;
    customerName: string;
    lastUpdTime: string;
    loginPassword: string;
    mobileNo: string;
    registerChannel: string;
    registerDate: string;
    status: string;
    customerId: string;
    avatar: string;
}
export type UpateUserInfoRequest = {
    verificationCode?: string,
    oldValue: string,
    updateValue: string,
    type: string
}
