// /store/userStore.ts
import { UserInfo } from "@/model/userModel";
import { makeAutoObservable } from 'mobx';

export class UserStore {
    accountType!: string;
    createTime!: string;
    customerId!: string;
    customerName!: string;
    lastupdTime!: string;
    loginDate!: string;
    loginPassword!: string;
    mobileNo!: string;
    registerChannel!: string;
    registerDate!: string;
    status!: string;
    // avatar: string = 'photograph';
    avatar: string = 'https://zhangyu-blog.oss-cn-beijing.aliyuncs.com/img/logo_20230516_uugai.com-1684220858267.png';
    protocol: boolean = false;

    clear() {
        this.accountType = ''
        this.createTime = ''
        this.customerId = ''
        this.lastupdTime = ''
        this.loginDate = ''
        this.loginPassword = ''
        this.mobileNo = ''
        this.registerChannel = ''
        this.registerDate = ''
        this.status = ''
        this.avatar = "https://zhangyu-blog.oss-cn-beijing.aliyuncs.com/img/logo_20230516_uugai.com-1684220858267.png"
    }

    setUserModel(user: UserInfo) {
        this.accountType = user.accountType
        this.createTime = user.registerDate
        this.customerName = user.customerName
        this.customerId = user.customerId
        this.lastupdTime = user.lastUpdTime
        this.loginDate = user.lastUpdTime
        this.mobileNo = user.mobileNo
        this.registerChannel = user.registerChannel
        this.registerDate = user.registerDate
        this.status = user.status
        this.avatar = user.avatar
    }

    constructor() {
        makeAutoObservable(this);
    }

    getAccountType(): string {
        return this.accountType;
    }

    setAccountType(value: string) {
        this.accountType = value;
    }

    getCreateTime(): string {
        return this.createTime;
    }

    setCreateTime(value: string) {
        this.createTime = value;
    }

    getCustomerId(): string {
        return this.customerId;
    }

    setCustomerId(value: string) {
        this.customerId = value;
    }

    getCustomerName(): string {
        return this.customerName;
    }

    setCustomerName(value: string) {
        this.customerName = value;
    }

    getLastupdTime(): string {
        return this.lastupdTime;
    }

    setLastupdTime(value: string) {
        this.lastupdTime = value;
    }

    getLoginDate(): string {
        return this.loginDate;
    }

    setLoginDate(value: string) {
        this.loginDate = value;
    }

    getLoginPassword(): string {
        return this.loginPassword;
    }

    setLoginPassword(value: string) {
        this.loginPassword = value;
    }

    getMobileNo(): string {
        return this.mobileNo;
    }

    setMobileNo(value: string) {
        this.mobileNo = value;
    }

    getRegisterChannel(): string {
        return this.registerChannel;
    }

    setRegisterChannel(value: string) {
        this.registerChannel = value;
    }

    getRegisterDate(): string {
        return this.registerDate;
    }

    setRegisterDate(value: string) {
        this.registerDate = value;
    }

    getStatus(): string {
        return this.status;
    }

    setStatus(value: string) {
        this.status = value;
    }

    getAvatar(): string {
        return this.avatar;
    }

    setAvatar(avatar: string) {
        this.avatar = avatar;
    }


}
