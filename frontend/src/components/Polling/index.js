import { Spin } from "antd";
import {
    LoadingOutlined
  } from '@ant-design/icons';
import { useEffect } from "react";
import { useRequest } from "ahooks";

export default function Polling({pollingFunc, fontSize}) {

    const { data, loading, run, cancel } = useRequest(pollingFunc, {
        pollingInterval: 3000,
        pollingWhenHidden: false,
    });
    
    useEffect(() => {
        run()
    },[])

    return (
        <Spin indicator={
            <LoadingOutlined
                style={{
                  fontSize: fontSize || 24,
                }}
                spin
            />} 
        />
    )
}