import { useEffect, useRef } from 'react';
import EChartsComponent from '../../components/EChartsComponent';

export default function HomePage(params) {
    
    return (
        <div className="home">
            <EChartsComponent 
                width="600px"
                height="400px"
                title="ECharts 入门示例"
                chartsOption={{
                    tooltip: {},
                    xAxis: {
                        data: ['衬衫', '羊毛衫', '雪纺衫', '裤子', '高跟鞋', '袜子']
                    },
                    yAxis: {},
                    series: [{
                        name: '销量',
                        type: 'bar',
                        data: [5, 20, 36, 10, 10, 20]
                    }]
                }}
            />
        </div>
    )
}