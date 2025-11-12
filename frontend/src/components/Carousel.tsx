import { cn } from "@/lib/utils";

interface Props extends React.HTMLAttributes<HTMLDivElement> {}

function Carousel({ className, children, ...props }: Props) {
    return (
        <div
            className={cn(
                "flex h-[22rem] w-full snap-x snap-proximity flex-row gap-4 overflow-x-scroll scroll-smooth px-2 pt-1 pb-3.5",
                className,
            )}
            tabIndex={-1}
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
        <div
            className={cn("h-full shrink-0 snap-start px-1", props)}
            {...props}
        >
            {children}
        </div>
    );
}

export { Carousel, CarouselItem };
