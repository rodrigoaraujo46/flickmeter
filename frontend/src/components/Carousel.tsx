import { cn } from "@/lib/utils";

interface Props extends React.HTMLAttributes<HTMLDivElement> {}

function Carousel({ className, children, ...props }: Props) {
    return (
        <div
            className={cn(
                "flex w-full snap-x snap-proximity flex-row gap-4 overflow-x-scroll scroll-smooth pb-3.5",
                className,
            )}
            style={{
                scrollbarColor: "grey transparent",
            }}
            {...props}
        >
            {children}
        </div>
    );
}

function CarouselItem({ className, children, ...props }: Props) {
    return (
        <div className={cn("h-60 shrink-0 snap-start", props)} {...props}>
            {children}
        </div>
    );
}

export { Carousel, CarouselItem };
