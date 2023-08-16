import { ipfsImg } from "../../request/api/public";
import { message, Upload } from "antd";

export const UploadProps = {
    name: "file",
    multiple: false,
    maxCount: 1,
    onStart(file) {
      console.log('onStart ==>', file);
    },
    onError(err) {
      console.log("onError ==>", err);
    },
    customRequest({
      data,
      file,
      filename,
      onError,
      onSuccess
    }) {
      // const isLt100M = file.size / 1024 / 1024 < 20;
      // if (!isLt100M) {
      //   return
      // }
      const formData = new FormData();
      if (data) {
        Object.keys(data).forEach(key => {
          formData.append(key, data[key]);
        });
      }
      formData.append('file', file);
      ipfsImg(formData)
      .then(res => {
        console.log('success ==>',res , file);
        onSuccess(res, file);
      })
      .catch(err => {
        console.log(err);
      });

      return {
        abort() {
          console.log("upload progress is aborted.");
        }
      };
    },
    beforeUpload(file) {
      const formatArr = ["image/jpeg","image/png","image/svg+xml","image/gif"]
      let isImage = false
      formatArr.map((e)=>{
        if ( file.type === e ) {
          isImage = true
        }
      })
      if (!isImage) {
        message.error("You can only upload JPG/PNG file!");
        return false || Upload.LIST_IGNORE
      }
      const isLt100M = file.size / 1024 / 1024 < 20;
      if (!isLt100M) {
        message.error("Image must smaller than 20MB!");
        return false || Upload.LIST_IGNORE
      }
    }
};