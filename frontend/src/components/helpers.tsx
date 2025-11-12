export const getRatingColor = (rating: number) => {
    if (rating <= 3) return "bg-red-600";
    if (rating <= 6) return "bg-yellow-500";
    if (rating <= 9) return "bg-green-600";
    return "bg-primary";
};
