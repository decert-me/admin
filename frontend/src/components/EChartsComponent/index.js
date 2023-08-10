import * as echarts from 'echarts';
import { useEffect, useRef } from "react"


export default function EChartsComponent(props) {
    
    const chartsRef = useRef();
    const { 
        title, 
        chartsOption,
        width,
        height
    } = props;

    function init() {
        const chart = echarts.init(chartsRef.current);
        chart.setOption({
            title: {
                text: title
            },
            ...chartsOption
        })
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div 
            className="chartsContainer" 
            ref={chartsRef}
            style={{
                width: width,
                height: height
            }}
        >

        </div>
    )
}