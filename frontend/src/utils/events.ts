import useShowMessage from "@/hooks/useShowMessage"
import { EventsOn } from "@/wailsjs/runtime/runtime"

const showMessage = useShowMessage()

//后端回调弹窗
EventsOn("globalMsg", async (re) => {
    if (re.type) {
        if (re.time && re.time > 0) {
            showMessage(re.type, re.msg, re.time)
        } else {
            showMessage(re.type, re.msg)
        }
    }
})
