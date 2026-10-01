import WishlistItem from "@/components/wishlist-item/wishlist-item";
import { AuthButton } from "@/components/auth-button/auth-button";
import { StylizedButton } from "@/components/ui/stylized-button";

export default function Page() {
  return (
      <main className="relative w-full max-w-6xl min-h-[700px] bg-[#FAF9F5]/90 backdrop-blur-sm border-[3px] border-[#214642] rounded-[40px] shadow-[0_20px_50px_rgba(33,70,66,0.1)] flex flex-col p-8 md:p-12 overflow-hidden">
        
        {/* Top Navigation Bar inside the rounded rectangle */}
        <header className="flex items-center justify-between w-full mb-12">
          {/* Top Left: Logo */}
          <div className="flex items-center">
            {/* <Logo /> */}
          </div>

          {/* Top Right: Auth Button */}
          <div className="flex items-center">
            <AuthButton />
          </div>
        </header>

        {/* Main Content Split */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-8 flex-grow items-center">
          
          {/* Left Column: WishlistItem Example */}
          <div className="flex justify-center items-center w-full">
            <div className="transform rotate-[-2deg] hover:rotate-0 transition-transform duration-300">
              <WishlistItem 
                name="Fujifilm X100VI"
                description="Compact digital camera, silver finish."
                price={1599}
                imageUrl="/assets/camera-sample.jpg"
                // comments={[
                //   { id: "1", author: "Alex", text: "I can pitch in $50!" },
                //   { id: "2", author: "Sam", text: "Great choice." }
                // ]}
              />
            </div>
          </div>

          {/* Right Column: Typography & Actions */}
          <div className="flex flex-col items-start justify-center space-y-10 max-w-md mx-auto lg:mx-0 lg:pl-8">
            
            {/* Header Text */}
            <h1 className="text-5xl md:text-6xl font-serif leading-[1.1] text-[#1E2423] tracking-tight">
              Let your friends <br />
              know what <br />
              <span className="italic text-[#214642]">you're wishing for!</span>
            </h1>

            {/* Action Buttons & Inputs */}
            <div className="w-full space-y-6">
              
              {/* Create Wishlist Button */}
              <StylizedButton
                variant="solid"
                size="lg"
                // onClick={handleCreateWishlist}
              >
                Create a Wishlist
              </StylizedButton>

              <div className="flex items-center gap-4 w-full">
                <div className="h-px bg-[#214642]/20 flex-grow"></div>
                <span className="text-[#214642]/60 text-sm font-medium uppercase tracking-widest">or</span>
                <div className="h-px bg-[#214642]/20 flex-grow"></div>
              </div>

              {/* Guest Share Token Entry */}
              <form 
              // onSubmit={handleViewWishlist} 
              className="flex flex-col space-y-3">
                <label htmlFor="token" className="text-sm font-semibold text-[#214642]">
                  View someone's wishlist
                </label>
                <div className="flex flex-col sm:flex-row gap-3">
                  <input
                    id="token"
                    type="text"
                    // value={shareToken}
                    // onChange={(e) => setShareToken(e.target.value)}
                    placeholder="Enter share token or URL..."
                    className="flex-grow py-3 px-4 rounded-xl bg-white border-2 border-[#214642]/40 text-[#1E2423] focus:outline-none focus:border-[#214642] transition-colors"
                  />
                  <StylizedButton type="submit" className="whitespace-nowrap">
                    View
                  </StylizedButton>
                </div>
              </form>
              
            </div>
          </div>
        </div>
      </main>
  );
}

