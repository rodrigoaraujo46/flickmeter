import type { Review } from "@/services/api/movies";
import { getRatingColor } from "./helpers";
import { ReviewDialog } from "./ReviewDialog";
import { Avatar, AvatarImage } from "./ui/avatar";
import { Button } from "./ui/button";

export function ReviewCard({ review }: { review: Review }) {
    return (
        <div
            key={review.id}
            className={`flex h-full w-full flex-col rounded-lg border border-border bg-card text-card-foreground`}
        >
            <div className="p-6">
                <div className="flex flex-row justify-between">
                    <div>
                        {review.updated_at && (
                            <p>
                                {new Date(review.updated_at).toLocaleDateString(
                                    undefined,
                                    {
                                        year: "numeric",
                                        month: "short",
                                        day: "numeric",
                                    },
                                )}
                            </p>
                        )}
                        <h3 className="mt-4 truncate font-bold text-2xl">
                            {review.title}
                        </h3>
                    </div>
                    <div
                        className={`${getRatingColor(review.rating)} flex aspect-square h-20 items-center justify-center rounded-lg text-center font-bold text-5xl text-primary-foreground`}
                    >
                        <p>{review.rating}</p>
                    </div>
                </div>
                <p className="wrap-break-word mt-1 line-clamp-5 whitespace-pre-wrap">
                    {review.review}
                </p>
            </div>
            <hr className="mt-auto border-border" />
            <div className="flex flex-row items-center justify-between p-2">
                <Button className="font-semibold" variant="link">
                    <Avatar className="w-auto">
                        <AvatarImage
                            referrerPolicy="no-referrer"
                            src={review.user.avatar_url}
                            alt={review.user.username}
                        />
                    </Avatar>
                    <p>{review.user.username}</p>
                </Button>
                <ReviewDialog review={review} />
            </div>
        </div>
    );
}
