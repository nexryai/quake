export const getQuakeScaleColor = (scale: number): string => {
    switch (scale) {
    case 10: return "#b3d5ff"; // 震度1 (薄い青)
    case 20: return "#00c980"; // 震度2 (青)
    case 30: return "#FFD700"; // 震度3 (水色)
    case 40: return "#FFA500"; // 震度4 (黄色)
    case 45: return "#f58b00"; // 震度5弱 (オレンジ)
    case 50: return "#ff5121"; // 震度5強 (濃いオレンジ)
    case 55: return "#FF0000"; // 震度6弱 (赤)
    case 60: return "#de0000"; // 震度6強 (濃い赤)
    case 70: return "#d700eb"; // 震度7 (紫)
    default: return "#CCCCCC"; // 不明な値 (グレー)
    }
};

export const getQuakeScaleColorText = (scale: number): string => {
    return scale <= 60 ? "white" : "black";
};